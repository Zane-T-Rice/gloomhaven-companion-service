package routers

import (
	"bytes"
	"context"
	"encoding/json"
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/middlewares"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofiber/fiber/v2"
)

// These need DynamoDB Local (scripts/start_dynamodb_local.sh) and are skipped without it.
func localDatabase(t *testing.T) utils.DynamoDB {
	t.Helper()
	os.Setenv(constants.LOCAL_SERVICE_PORT, "7575")
	os.Setenv(constants.LOCAL_DATABASE_ENDPOINT, "http://localhost:8000/")
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	db := utils.NewDynamoDB()
	db.ConnectToDynamoDB()
	_, err := db.DynamoDBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(constants.TABLE_NAME),
		AttributeDefinitions: []dbtypes.AttributeDefinition{
			{AttributeName: aws.String("parent"), AttributeType: dbtypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("entity"), AttributeType: dbtypes.ScalarAttributeTypeS},
		},
		KeySchema: []dbtypes.KeySchemaElement{
			{AttributeName: aws.String("parent"), KeyType: dbtypes.KeyTypeHash},
			{AttributeName: aws.String("entity"), KeyType: dbtypes.KeyTypeRange},
		},
		BillingMode: dbtypes.BillingModePayPerRequest,
		GlobalSecondaryIndexes: []dbtypes.GlobalSecondaryIndex{{
			IndexName: aws.String("entity-index"),
			KeySchema: []dbtypes.KeySchemaElement{
				{AttributeName: aws.String("entity"), KeyType: dbtypes.KeyTypeHash},
				{AttributeName: aws.String("parent"), KeyType: dbtypes.KeyTypeRange},
			},
			Projection: &dbtypes.Projection{ProjectionType: dbtypes.ProjectionTypeKeysOnly},
		}},
	})
	// The table already existing is fine. Anything else means there is no DynamoDB Local.
	if err != nil && db.Ping() != nil {
		t.Skip("DynamoDB Local is not running")
	}
	return db
}

type caller struct {
	app   *fiber.App
	id    string
	admin bool
}

// Acts as the given player, without a real token.
func (c caller) call(t *testing.T, method string, path string, body any) (int, []byte) {
	t.Helper()
	payload, _ := json.Marshal(body)
	request := httptest.NewRequest(method, path, bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Player", c.id)
	if c.admin {
		request.Header.Set("X-Admin", "yes")
	}
	response, err := c.app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(response.Body)
	return response.StatusCode, raw
}

func newApp(db utils.DynamoDB) *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		scope := constants.SCOPE_PUBLIC
		if c.Get("X-Admin") != "" {
			scope = constants.SCOPE_ADMIN
		}
		c.Context().SetUserValue(jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
			RegisteredClaims: validator.RegisteredClaims{Subject: c.Get("X-Player")},
			CustomClaims:     &middlewares.CustomClaims{Scope: scope},
		})
		return c.Next()
	})
	RegisterTemplatesRoutes(app, db)
	return app
}

type templateBody struct {
	Entity    string `json:"entity"`
	Parent    string `json:"parent"`
	Type      string `json:"type"`
	Scope     string `json:"scope"`
	EditionId string `json:"editionId"`
	UpdatedAt string `json:"updatedAt"`
}

func ids(t *testing.T, raw []byte) map[string]templateBody {
	t.Helper()
	list := []templateBody{}
	if err := json.Unmarshal(raw, &list); err != nil {
		t.Fatalf("%s: %v", raw, err)
	}
	byType := map[string]templateBody{}
	for _, item := range list {
		byType[item.Type] = item
	}
	return byType
}

func TestTemplatesAndEditionsBelongToTheirOwners(t *testing.T) {
	db := localDatabase(t)
	app := newApp(db)
	suffix := t.Name()
	admin := caller{app, "admin|" + suffix, true}
	alice := caller{app, "alice|" + suffix, false}
	bob := caller{app, "bob|" + suffix, false}
	carol := caller{app, "carol|" + suffix, false}

	stats := map[string]any{"0": map[string]any{"normal": map[string]any{"maximumHP": 3, "class": "x"}}}
	system := func(name string) templateBody {
		_, raw := admin.call(t, "POST", "/templates?scope=system", map[string]any{"type": name, "stats": stats})
		body := templateBody{}
		json.Unmarshal(raw, &body)
		return body
	}
	sysTemplate := system("Old System Thing " + suffix)

	// A template without an edition is in the default one, and is the system's.
	if sysTemplate.Scope != "system" || sysTemplate.EditionId != constants.DEFAULT_EDITION_ID {
		t.Fatalf("system template: %+v", sysTemplate)
	}

	// Only admins make system templates.
	if status, _ := alice.call(t, "POST", "/templates?scope=system", map[string]any{"type": "No", "stats": stats}); status == 200 {
		t.Fatal("a player made a system template")
	}

	// Alice makes an edition and a template in it.
	_, raw := alice.call(t, "POST", "/templates/editions", map[string]any{"name": "Homebrew"})
	edition := types.EditionItem{}
	json.Unmarshal(raw, &edition)
	editionId := edition.Entity[len("#EDITION#"):]
	status, raw := alice.call(t, "POST", "/templates", map[string]any{"type": "Alice Imp " + suffix, "stats": stats, "editionId": editionId})
	mine := templateBody{}
	json.Unmarshal(raw, &mine)
	if status != 200 || mine.Scope != "mine" || mine.EditionId != editionId {
		t.Fatalf("alice template: %d %s", status, raw)
	}

	// Bob is in no campaign with Alice, so he sees the system's templates but not hers.
	_, raw = bob.call(t, "GET", "/templates", nil)
	seen := ids(t, raw)
	if _, ok := seen[sysTemplate.Type]; !ok {
		t.Fatal("bob can not see the system template")
	}
	if _, ok := seen[mine.Type]; ok {
		t.Fatal("bob can see alice's template")
	}

	// Once they share a campaign, he can see it, read only, and without her login id in it.
	for _, player := range []caller{alice, bob} {
		db.PutItem(types.PlayerCampaignItem{Item: types.Item{
			Parent: constants.PLAYER + "#" + player.id,
			Entity: constants.CAMPAIGN + "#shared-" + suffix,
		}})
	}
	_, raw = bob.call(t, "GET", "/templates", nil)
	shared, ok := ids(t, raw)[mine.Type]
	if !ok || shared.Scope != "shared" || shared.Parent != constants.PLAYER {
		t.Fatalf("bob sees %+v", shared)
	}
	_, raw = bob.call(t, "GET", "/templates/editions", nil)
	if !bytes.Contains(raw, []byte("Homebrew")) {
		t.Fatalf("bob can not see the edition: %s", raw)
	}
	if status, _ := bob.call(t, "PATCH", "/templates/"+mine.Entity[len("#TEMPLATE#"):], map[string]any{"type": "Bob's now", "updatedAt": mine.UpdatedAt}); status == 200 {
		t.Fatal("bob changed alice's template")
	}
	if status, _ := bob.call(t, "DELETE", "/templates/"+mine.Entity[len("#TEMPLATE#"):], nil); status == 200 {
		t.Fatal("bob deleted alice's template")
	}

	// Carol is in no campaign with Alice, so she sees nothing of hers.
	_, raw = carol.call(t, "GET", "/templates", nil)
	if _, ok := ids(t, raw)[mine.Type]; ok {
		t.Fatal("carol can see alice's template")
	}

	// A player can not change the system's templates, an admin can.
	sysId := sysTemplate.Entity[len("#TEMPLATE#"):]
	if status, _ := alice.call(t, "PATCH", "/templates/"+sysId, map[string]any{"standeeLimit": 3, "updatedAt": sysTemplate.UpdatedAt}); status == 200 {
		t.Fatal("alice changed a system template")
	}
	if status, raw := admin.call(t, "PATCH", "/templates/"+sysId, map[string]any{"standeeLimit": 3, "updatedAt": sysTemplate.UpdatedAt}); status != 200 {
		t.Fatalf("admin could not change a system template: %s", raw)
	}
	if status, _ := alice.call(t, "DELETE", "/templates/"+sysId, nil); status == 200 {
		t.Fatal("alice deleted a system template")
	}

	// Bob can not put a template in Alice's edition.
	if status, _ := bob.call(t, "POST", "/templates", map[string]any{"type": "Bob Imp " + suffix, "stats": stats, "editionId": editionId}); status == 200 {
		t.Fatal("bob used alice's edition")
	}
	// A system template can not go in her edition either.
	if status, _ := admin.call(t, "POST", "/templates?scope=system", map[string]any{"type": "Sys In Hers " + suffix, "stats": stats, "editionId": editionId}); status == 200 {
		t.Fatal("a system template went in a player's edition")
	}

	// An edition with templates in it is kept. Alice can move her template out and then delete it.
	if status, _ := alice.call(t, "DELETE", "/templates/editions/"+editionId, nil); status == 200 {
		t.Fatal("deleted an edition that still has templates")
	}
	mineId := mine.Entity[len("#TEMPLATE#"):]
	if status, raw := alice.call(t, "DELETE", "/templates/"+mineId, nil); status != 200 {
		t.Fatalf("alice could not delete her template: %s", raw)
	}
	if status, raw := alice.call(t, "DELETE", "/templates/editions/"+editionId, nil); status != 200 {
		t.Fatalf("alice could not delete her empty edition: %s", raw)
	}

	// The default edition is always there and can not be removed.
	_, raw = bob.call(t, "GET", "/templates/editions", nil)
	if !bytes.Contains(raw, []byte(constants.DEFAULT_EDITION_NAME)) {
		t.Fatalf("no default edition: %s", raw)
	}
	if status, _ := admin.call(t, "DELETE", "/templates/editions/"+constants.DEFAULT_EDITION_ID, nil); status == 200 {
		t.Fatal("deleted the default edition")
	}

	admin.call(t, "DELETE", "/templates/"+sysId, nil)
}
