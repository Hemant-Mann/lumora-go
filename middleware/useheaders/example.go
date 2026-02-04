package useheaders

// Example usage of UseHeaders middleware:
//
// import (
//     z "github.com/Oudwins/zog"
//     "github.com/hemant-mann/lumora-go/core"
//     "github.com/hemant-mann/lumora-go/middleware/useheaders"
// )
//
// // Define your schema using zog
// // Note: Header names are normalized to lowercase for schema matching
// var authHeadersSchema = z.Struct(z.Shape{
//     "authorization": z.String().Min(1),
//     "x-api-key":     z.String().Optional(),
//     "content-type":  z.String().Optional(),
// })
//
// // Define your struct
// type AuthHeaders struct {
//     Authorization string `json:"authorization"`
//     APIKey        string `json:"x-api-key,omitempty"`
//     ContentType   string `json:"content-type,omitempty"`
// }
//
// // Use in route
// app.Get("/protected",
//     func(ctx core.Context) (*core.Response, error) {
//         // Get parsed and validated headers
//         headers := useheaders.GetHeaders(ctx).(*AuthHeaders)
//
//         // Use validated headers
//         resp := core.NewResponse().
//             WithStatus(200).
//             WithBody(map[string]any{
//                 "message": "Access granted",
//                 "token":   headers.Authorization,
//             })
//         return resp, nil
//     },
//     useheaders.UseHeaders(authHeadersSchema, &AuthHeaders{}),
// )
