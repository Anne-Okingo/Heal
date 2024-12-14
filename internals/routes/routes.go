package routes

// Allowed routes
var allowedRoutes = map[string]bool{
	"/":                 true,
	"/privacy":          true,
	"/session":          true,
	"/signup":           true,
	"/loggin":           true,
	"/welcome":          true,
	"/api/get-username": true,
	"/api/logout":       true,
	"/api/gemini":       true,
	"/api/speechify":    true,
	"/get-heard": true,
}
