package handlers

import (
	"encoding/json"
)

// GenerateOpenAPISpec generates OpenAPI 3.0 spec from registered routes
func GenerateOpenAPISpec(baseURL string) string {
	spec := make(map[string]interface{})
	spec["openapi"] = "3.0.3"
	
	info := make(map[string]interface{})
	info["title"] = "Go Produk API"
	info["version"] = "1.0.0"
	info["description"] = "Interactive API documentation for Go Produk service"
	spec["info"] = info
	
	servers := make([]interface{}, 1)
	servers[0] = map[string]string{"url": baseURL}
	spec["servers"] = servers
	
	spec["paths"] = generatePaths()
	spec["components"] = generateComponents()
	
	data, _ := json.MarshalIndent(spec, "", "  ")
	return string(data)
}

func generatePaths() map[string]interface{} {
	paths := make(map[string]interface{})
	
	// Health endpoint
	healthGet := make(map[string]interface{})
	healthGet["summary"] = "Health check"
	healthGet["description"] = "Check health of all databases"
	healthResp := make(map[string]interface{})
	healthResp["200"] = map[string]string{"description": "Healthy"}
	healthGet["responses"] = healthResp
	paths["/health"] = map[string]interface{}{"get": healthGet}
	
	// Auth register
	authRegPost := make(map[string]interface{})
	authRegPost["summary"] = "Register new user"
	authRegPost["description"] = "Register a new user account"
	authRegBody := make(map[string]interface{})
	authRegBody["required"] = true
	authRegContent := make(map[string]interface{})
	authRegJSON := make(map[string]interface{})
	authRegJSON["schema"] = map[string]interface{}{"$ref": "#/components/schemas/RegisterRequest"}
	authRegContent["application/json"] = authRegJSON
	authRegBody["content"] = authRegContent
	authRegPost["requestBody"] = authRegBody
	authRegResp := make(map[string]interface{})
	authRegResp["201"] = map[string]string{"description": "User created"}
	authRegPost["responses"] = authRegResp
	paths["/api/v1/auth/register"] = map[string]interface{}{"post": authRegPost}
	
	// Auth login
	authLoginPost := make(map[string]interface{})
	authLoginPost["summary"] = "Login user"
	authLoginPost["description"] = "Authenticate user with email and password"
	loginBody := make(map[string]interface{})
	loginBody["required"] = true
	loginContent := make(map[string]interface{})
	loginJSON := make(map[string]interface{})
	loginJSON["schema"] = map[string]interface{}{"$ref": "#/components/schemas/LoginRequest"}
	loginContent["application/json"] = loginJSON
	loginBody["content"] = loginContent
	authLoginPost["requestBody"] = loginBody
	loginResp := make(map[string]interface{})
	loginResp["200"] = map[string]string{"description": "Login successful"}
	authLoginPost["responses"] = loginResp
	paths["/api/v1/auth/login"] = map[string]interface{}{"post": authLoginPost}
	
	// Auth logout
	authLogoutPost := make(map[string]interface{})
	authLogoutPost["summary"] = "Logout user"
	authLogoutPost["description"] = "Logout and invalidate token"
	authLogoutSec := make([]interface{}, 1)
	authLogoutSec[0] = map[string][]string{"BearerAuth": {}}
	authLogoutPost["security"] = authLogoutSec
	authLogoutResp := make(map[string]interface{})
	authLogoutResp["200"] = map[string]string{"description": "Logged out"}
	authLogoutPost["responses"] = authLogoutResp
	paths["/api/v1/auth/logout"] = map[string]interface{}{"post": authLogoutPost}
	
	// Users /me - GET
	usersMeGet := make(map[string]interface{})
	usersMeGet["summary"] = "Get current user profile"
	usersMeGet["description"] = "Get profile of authenticated user"
	usersMeSec := make([]interface{}, 1)
	usersMeSec[0] = map[string][]string{"BearerAuth": {}}
	usersMeGet["security"] = usersMeSec
	usersMeResp := make(map[string]interface{})
	usersMeResp["200"] = map[string]interface{}{
		"description": "User profile",
		"content": map[string]interface{}{
			"application/json": map[string]interface{}{
				"schema": map[string]interface{}{"$ref": "#/components/schemas/User"},
			},
		},
	}
	usersMeGet["responses"] = usersMeResp
	usersMe := make(map[string]interface{})
	usersMe["get"] = usersMeGet
	paths["/api/v1/users/me"] = usersMe
	
	// Users /me - PUT
	usersMePut := make(map[string]interface{})
	usersMePut["summary"] = "Update user profile"
	usersMePut["description"] = "Update current user profile"
	usersMePut["security"] = usersMeSec
	usersMePutBody := make(map[string]interface{})
	usersMePutBody["required"] = true
	usersMePutContent := make(map[string]interface{})
	usersMePutJSON := make(map[string]interface{})
	usersMePutJSON["schema"] = map[string]interface{}{"$ref": "#/components/schemas/UpdateProfileRequest"}
	usersMePutContent["application/json"] = usersMePutJSON
	usersMePutBody["content"] = usersMePutContent
	usersMePut["requestBody"] = usersMePutBody
	usersMePutResp := make(map[string]interface{})
	usersMePutResp["200"] = map[string]string{"description": "Profile updated"}
	usersMePut["responses"] = usersMePutResp
	usersMe["put"] = usersMePut
	
	// Categories - GET all
	categoriesGet := make(map[string]interface{})
	categoriesGet["summary"] = "Get all categories"
	categoriesGet["description"] = "List all categories (public)"
	categoriesParams := make([]interface{}, 1)
	categoriesParams[0] = map[string]interface{}{"name": "tree", "in": "query", "schema": map[string]interface{}{"type": "boolean"}}
	categoriesGet["parameters"] = categoriesParams
	categoriesResp := make(map[string]interface{})
	categoriesResp["200"] = map[string]string{"description": "List of categories"}
	categoriesGet["responses"] = categoriesResp
	paths["/api/v1/categories"] = map[string]interface{}{"get": categoriesGet}
	
	// Categories - POST create
	categoriesPost := make(map[string]interface{})
	categoriesPost["summary"] = "Create category (Admin only)"
	categoriesPost["description"] = "Create a new category"
	categoriesPost["security"] = []interface{}{map[string][]string{"BearerAuth": {}}}
	categoriesPostBody := make(map[string]interface{})
	categoriesPostBody["required"] = true
	categoriesPostContent := make(map[string]interface{})
	categoriesPostJSON := make(map[string]interface{})
	categoriesPostJSON["schema"] = map[string]interface{}{"$ref": "#/components/schemas/CreateCategoryRequest"}
	categoriesPostContent["application/json"] = categoriesPostJSON
	categoriesPostBody["content"] = categoriesPostContent
	categoriesPost["requestBody"] = categoriesPostBody
	categoriesPostResp := make(map[string]interface{})
	categoriesPostResp["201"] = map[string]string{"description": "Category created"}
	categoriesPost["responses"] = categoriesPostResp
	paths["/api/v1/categories"].(map[string]interface{})["post"] = categoriesPost
	
	// Categories/{id} - GET
	catIdGet := make(map[string]interface{})
	catIdGet["summary"] = "Get category by ID"
	catIdGet["description"] = "Get single category by ID"
	catIdParams := make([]interface{}, 1)
	catIdParams[0] = map[string]interface{}{"name": "id", "in": "path", "required": true, "schema": map[string]interface{}{"type": "string"}}
	catIdGet["parameters"] = catIdParams
	catIdResp := make(map[string]interface{})
	catIdResp["200"] = map[string]string{"description": "Category details"}
	catIdGet["responses"] = catIdResp
	paths["/api/v1/categories/{id}"] = map[string]interface{}{"get": catIdGet}
	
	// Categories/{id} - PUT
	catIdPut := make(map[string]interface{})
	catIdPut["summary"] = "Update category (Admin only)"
	catIdPut["description"] = "Update category by ID"
	catIdPut["security"] = []interface{}{map[string][]string{"BearerAuth": {}}}
	catIdPut["parameters"] = catIdParams
	catIdPutBody := make(map[string]interface{})
	catIdPutBody["required"] = true
	catIdPutContent := make(map[string]interface{})
	catIdPutJSON := make(map[string]interface{})
	catIdPutJSON["schema"] = map[string]interface{}{"$ref": "#/components/schemas/UpdateCategoryRequest"}
	catIdPutContent["application/json"] = catIdPutJSON
	catIdPutBody["content"] = catIdPutContent
	catIdPut["requestBody"] = catIdPutBody
	catIdPutResp := make(map[string]interface{})
	catIdPutResp["200"] = map[string]string{"description": "Category updated"}
	catIdPut["responses"] = catIdPutResp
	paths["/api/v1/categories/{id}"].(map[string]interface{})["put"] = catIdPut
	
	// Categories/{id} - DELETE
	catIdDel := make(map[string]interface{})
	catIdDel["summary"] = "Delete category (Admin only)"
	catIdDel["description"] = "Delete category by ID"
	catIdDel["security"] = []interface{}{map[string][]string{"BearerAuth": {}}}
	catIdDel["parameters"] = catIdParams
	catIdDelResp := make(map[string]interface{})
	catIdDelResp["204"] = map[string]string{"description": "Category deleted"}
	catIdDel["responses"] = catIdDelResp
	paths["/api/v1/categories/{id}"].(map[string]interface{})["delete"] = catIdDel
	
	// Products - GET all
	productsGet := make(map[string]interface{})
	productsGet["summary"] = "Get all products"
	productsGet["description"] = "List all products (public)"
	productsParams := make([]interface{}, 3)
	productsParams[0] = map[string]interface{}{"name": "category_id", "in": "query", "schema": map[string]interface{}{"type": "string"}}
	productsParams[1] = map[string]interface{}{"name": "page", "in": "query", "schema": map[string]interface{}{"type": "integer", "default": 1}}
	productsParams[2] = map[string]interface{}{"name": "limit", "in": "query", "schema": map[string]interface{}{"type": "integer", "default": 20}}
	productsGet["parameters"] = productsParams
	productsResp := make(map[string]interface{})
	productsResp["200"] = map[string]string{"description": "List of products"}
	productsGet["responses"] = productsResp
	paths["/api/v1/products"] = map[string]interface{}{"get": productsGet}
	
	// Products - POST create
	productsPost := make(map[string]interface{})
	productsPost["summary"] = "Create product (Admin only)"
	productsPost["description"] = "Create a new product"
	productsPost["security"] = []interface{}{map[string][]string{"BearerAuth": {}}}
	productsPostBody := make(map[string]interface{})
	productsPostBody["required"] = true
	productsPostContent := make(map[string]interface{})
	productsPostJSON := make(map[string]interface{})
	productsPostJSON["schema"] = map[string]interface{}{"$ref": "#/components/schemas/CreateProductRequest"}
	productsPostContent["application/json"] = productsPostJSON
	productsPostBody["content"] = productsPostContent
	productsPost["requestBody"] = productsPostBody
	productsPostResp := make(map[string]interface{})
	productsPostResp["201"] = map[string]string{"description": "Product created"}
	productsPost["responses"] = productsPostResp
	paths["/api/v1/products"].(map[string]interface{})["post"] = productsPost
	
	// Products/{id} - GET
	prodIdGet := make(map[string]interface{})
	prodIdGet["summary"] = "Get product by ID"
	prodIdGet["description"] = "Get single product by ID"
	prodIdParams := make([]interface{}, 1)
	prodIdParams[0] = map[string]interface{}{"name": "id", "in": "path", "required": true, "schema": map[string]interface{}{"type": "string"}}
	prodIdGet["parameters"] = prodIdParams
	prodIdResp := make(map[string]interface{})
	prodIdResp["200"] = map[string]string{"description": "Product details"}
	prodIdGet["responses"] = prodIdResp
	paths["/api/v1/products/{id}"] = map[string]interface{}{"get": prodIdGet}
	
	// Products/{id} - PUT
	prodIdPut := make(map[string]interface{})
	prodIdPut["summary"] = "Update product (Admin only)"
	prodIdPut["description"] = "Update product by ID"
	prodIdPut["security"] = []interface{}{map[string][]string{"BearerAuth": {}}}
	prodIdPut["parameters"] = prodIdParams
	prodIdPutBody := make(map[string]interface{})
	prodIdPutBody["required"] = true
	prodIdPutContent := make(map[string]interface{})
	prodIdPutJSON := make(map[string]interface{})
	prodIdPutJSON["schema"] = map[string]interface{}{"$ref": "#/components/schemas/UpdateProductRequest"}
	prodIdPutContent["application/json"] = prodIdPutJSON
	prodIdPutBody["content"] = prodIdPutContent
	prodIdPut["requestBody"] = prodIdPutBody
	prodIdPutResp := make(map[string]interface{})
	prodIdPutResp["200"] = map[string]string{"description": "Product updated"}
	prodIdPut["responses"] = prodIdPutResp
	paths["/api/v1/products/{id}"].(map[string]interface{})["put"] = prodIdPut
	
	// Products/{id} - DELETE
	prodIdDel := make(map[string]interface{})
	prodIdDel["summary"] = "Delete product (Admin only)"
	prodIdDel["description"] = "Delete product by ID"
	prodIdDel["security"] = []interface{}{map[string][]string{"BearerAuth": {}}}
	prodIdDel["parameters"] = prodIdParams
	prodIdDelResp := make(map[string]interface{})
	prodIdDelResp["204"] = map[string]string{"description": "Product deleted"}
	prodIdDel["responses"] = prodIdDelResp
	paths["/api/v1/products/{id}"].(map[string]interface{})["delete"] = prodIdDel
	
	return paths
}

func generateComponents() map[string]interface{} {
	components := make(map[string]interface{})
	schemas := make(map[string]interface{})
	
	// LoginRequest
	loginReq := make(map[string]interface{})
	loginReq["type"] = "object"
	loginProps := make(map[string]interface{})
	loginProps["email"] = map[string]interface{}{"type": "string", "format": "email"}
	loginProps["password"] = map[string]interface{}{"type": "string"}
	loginReq["properties"] = loginProps
	loginReq["required"] = []interface{}{"email", "password"}
	schemas["LoginRequest"] = loginReq
	
	// RegisterRequest
	regReq := make(map[string]interface{})
	regReq["type"] = "object"
	regProps := make(map[string]interface{})
	regProps["name"] = map[string]interface{}{"type": "string"}
	regProps["email"] = map[string]interface{}{"type": "string", "format": "email"}
	regProps["password"] = map[string]interface{}{"type": "string", "minLength": 6}
	regReq["properties"] = regProps
	regReq["required"] = []interface{}{"name", "email", "password"}
	schemas["RegisterRequest"] = regReq
	
	// LoginResponse
	loginResp := make(map[string]interface{})
	loginResp["type"] = "object"
	loginRespProps := make(map[string]interface{})
	loginRespProps["user"] = map[string]interface{}{"$ref": "#/components/schemas/User"}
	loginRespProps["access_token"] = map[string]interface{}{"type": "string"}
	loginRespProps["expires_in"] = map[string]interface{}{"type": "integer"}
	loginResp["properties"] = loginRespProps
	schemas["LoginResponse"] = loginResp
	
	// User
	user := make(map[string]interface{})
	user["type"] = "object"
	userProps := make(map[string]interface{})
	userProps["id"] = map[string]interface{}{"type": "string", "format": "uuid"}
	userProps["name"] = map[string]interface{}{"type": "string"}
	userProps["email"] = map[string]interface{}{"type": "string", "format": "email"}
	userProps["role"] = map[string]interface{}{"type": "string", "enum": []interface{}{"user", "admin"}}
	userProps["is_active"] = map[string]interface{}{"type": "boolean"}
	userProps["created_at"] = map[string]interface{}{"type": "string", "format": "date-time"}
	userProps["updated_at"] = map[string]interface{}{"type": "string", "format": "date-time"}
	user["properties"] = userProps
	schemas["User"] = user
	
	// UpdateProfileRequest
	updateProfile := make(map[string]interface{})
	updateProfile["type"] = "object"
	updateProfileProps := make(map[string]interface{})
	updateProfileProps["name"] = map[string]interface{}{"type": "string"}
	updateProfileProps["username"] = map[string]interface{}{"type": "string"}
	updateProfileProps["display_name"] = map[string]interface{}{"type": "string"}
	updateProfileProps["phone"] = map[string]interface{}{"type": "string"}
	updateProfileProps["bio"] = map[string]interface{}{"type": "string"}
	updateProfileProps["avatar_url"] = map[string]interface{}{"type": "string", "format": "uri"}
	updateProfileProps["location"] = map[string]interface{}{"type": "string"}
	updateProfile["properties"] = updateProfileProps
	schemas["UpdateProfileRequest"] = updateProfile
	
	// CreateCategoryRequest
	createCat := make(map[string]interface{})
	createCat["type"] = "object"
	createCatProps := make(map[string]interface{})
	createCatProps["name"] = map[string]interface{}{"type": "string"}
	createCatProps["slug"] = map[string]interface{}{"type": "string"}
	createCatProps["icon"] = map[string]interface{}{"type": "string"}
	createCatProps["description"] = map[string]interface{}{"type": "string"}
	createCatProps["sort_order"] = map[string]interface{}{"type": "integer"}
	createCatProps["is_active"] = map[string]interface{}{"type": "boolean"}
	createCat["properties"] = createCatProps
	createCat["required"] = []interface{}{"name", "slug"}
	schemas["CreateCategoryRequest"] = createCat
	
	// UpdateCategoryRequest
	updateCat := make(map[string]interface{})
	updateCat["type"] = "object"
	updateCatProps := make(map[string]interface{})
	updateCatProps["name"] = map[string]interface{}{"type": "string"}
	updateCatProps["slug"] = map[string]interface{}{"type": "string"}
	updateCatProps["icon"] = map[string]interface{}{"type": "string"}
	updateCatProps["description"] = map[string]interface{}{"type": "string"}
	updateCatProps["sort_order"] = map[string]interface{}{"type": "integer"}
	updateCatProps["is_active"] = map[string]interface{}{"type": "boolean"}
	updateCat["properties"] = updateCatProps
	schemas["UpdateCategoryRequest"] = updateCat
	
	// CreateProductRequest
	createProd := make(map[string]interface{})
	createProd["type"] = "object"
	createProdProps := make(map[string]interface{})
	createProdProps["category_id"] = map[string]interface{}{"type": "string", "format": "uuid"}
	createProdProps["name"] = map[string]interface{}{"type": "string"}
	createProdProps["slug"] = map[string]interface{}{"type": "string"}
	createProdProps["description"] = map[string]interface{}{"type": "string"}
	createProdProps["price"] = map[string]interface{}{"type": "integer"}
	createProdProps["stock"] = map[string]interface{}{"type": "integer"}
	createProdProps["is_active"] = map[string]interface{}{"type": "boolean"}
	createProd["properties"] = createProdProps
	createProd["required"] = []interface{}{"category_id", "name", "slug", "price"}
	schemas["CreateProductRequest"] = createProd
	
	// UpdateProductRequest
	updateProd := make(map[string]interface{})
	updateProd["type"] = "object"
	updateProdProps := make(map[string]interface{})
	updateProdProps["category_id"] = map[string]interface{}{"type": "string", "format": "uuid"}
	updateProdProps["name"] = map[string]interface{}{"type": "string"}
	updateProdProps["slug"] = map[string]interface{}{"type": "string"}
	updateProdProps["description"] = map[string]interface{}{"type": "string"}
	updateProdProps["price"] = map[string]interface{}{"type": "integer"}
	updateProdProps["stock"] = map[string]interface{}{"type": "integer"}
	updateProdProps["is_active"] = map[string]interface{}{"type": "boolean"}
	updateProd["properties"] = updateProdProps
	schemas["UpdateProductRequest"] = updateProd
	
	components["schemas"] = schemas
	
	// Security schemes
	securitySchemes := make(map[string]interface{})
	bearerAuth := make(map[string]interface{})
	bearerAuth["type"] = "http"
	bearerAuth["scheme"] = "bearer"
	bearerAuth["bearerFormat"] = "JWT"
	securitySchemes["BearerAuth"] = bearerAuth
	components["securitySchemes"] = securitySchemes
	
	return components
}
