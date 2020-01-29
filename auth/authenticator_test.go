package auth

import (
	"encoding/base64"
	"net/http/httptest"
	"testing"

	"github.com/System-Glitch/goyave/v2"
	"github.com/System-Glitch/goyave/v2/config"
	"github.com/System-Glitch/goyave/v2/database"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type TestUser struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100)" auth:"password"`
	Email    string `gorm:"type:varchar(100);unique_index" auth:"username"`
}

type TestUserPromoted struct {
	TestUser
}

type TestUserOverride struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100);column:password_override" auth:"password"`
	Email    string `gorm:"type:varchar(100);unique_index" auth:"username"`
}

type AuthenticationTestSuite struct {
	goyave.TestSuite
}

func (suite *AuthenticationTestSuite) SetupSuite() {
	config.Set("dbConnection", "mysql")
	database.ClearRegisteredModels()
	database.RegisterModel(&TestUser{})

	database.Migrate()
}

func (suite *AuthenticationTestSuite) SetupTest() {
	user := &TestUser{
		Name:     "Admin",
		Password: "$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password"
		Email:    "johndoe@example.org",
	}
	database.GetConnection().Create(user)
}

func (suite *AuthenticationTestSuite) TestFindColumns() {
	user := &TestUser{}
	fields := FindColumns(user, "username", "password")
	suite.Len(fields, 2)
	suite.Equal("email", fields[0].Name)
	suite.Equal("password", fields[1].Name)

	fields = FindColumns(user, "username", "notatag", "password")
	suite.Len(fields, 3)
	suite.Equal("email", fields[0].Name)
	suite.Nil(fields[1])
	suite.Equal("password", fields[2].Name)

	userOverride := &TestUserOverride{}
	fields = FindColumns(userOverride, "password")
	suite.Len(fields, 1)
	suite.Equal("password_override", fields[0].Name)
}

func (suite *AuthenticationTestSuite) TestFindColumnsPromoted() {
	user := &TestUserPromoted{}
	fields := FindColumns(user, "username", "password")
	suite.Len(fields, 2)
	suite.Equal("email", fields[0].Name)
	suite.Equal("password", fields[1].Name)

	fields = FindColumns(user, "username", "notatag", "password")
	suite.Len(fields, 3)
	suite.Equal("email", fields[0].Name)
	suite.Nil(fields[1])
	suite.Equal("password", fields[2].Name)
}

func (suite *AuthenticationTestSuite) TestAuthMiddleware() {
	// Test middleware with BasicAuth
	authenticator := Middleware(&TestUser{}, &BasicAuthenticator{})

	request := suite.CreateTestRequest(httptest.NewRequest("GET", "/", nil))
	request.Header().Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("johndoe@example.org:wrong_password")))
	result := suite.Middleware(authenticator, request, func(response *goyave.Response, request *goyave.Request) {
		suite.Fail("Auth middleware passed")
	})
	suite.Equal(401, result.StatusCode)

	request = suite.CreateTestRequest(httptest.NewRequest("GET", "/", nil))
	request.Header().Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("johndoe@example.org:password")))
	result = suite.Middleware(authenticator, request, func(response *goyave.Response, request *goyave.Request) {
		suite.IsType(&TestUser{}, request.User)
		suite.Equal("Admin", request.User.(*TestUser).Name)
		response.Status(200)
	})
	suite.Equal(200, result.StatusCode)
}

func (suite *AuthenticationTestSuite) TestGetBearerToken() {
	request := suite.CreateTestRequest(nil)
	request.Header().Set("Authorization", "NotBearer 123456789")
	suite.Empty(GetBearerToken(request))

	request.Header().Set("Authorization", "Bearer123456789")
	suite.Empty(GetBearerToken(request))

	request.Header().Set("Authorization", "Bearer 123456789")
	suite.Equal("123456789", GetBearerToken(request))
}

func (suite *AuthenticationTestSuite) TearDownTest() {
	suite.ClearDatabase()
}

func (suite *AuthenticationTestSuite) TearDownSuite() {
	database.GetConnection().DropTable(&TestUser{})
	database.ClearRegisteredModels()
}

func TestAuthenticationSuite(t *testing.T) {
	goyave.RunTest(t, new(AuthenticationTestSuite))
}
