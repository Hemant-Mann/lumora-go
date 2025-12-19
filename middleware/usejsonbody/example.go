package usejsonbody

// Example usage of UseJsonBody middleware:
//
// import (
//     z "github.com/Oudwins/zog"
//     "github.com/hemant-mann/lumora-go/core"
//     "github.com/hemant-mann/lumora-go/middleware/usejsonbody"
// )
//
// // Define your schema using zog (returns *z.StructSchema which implements SchemaWithParse)
// var userSchema = z.Struct(z.Shape{
//     "name":  z.String().Min(3).Max(50),
//     "email": z.String().Email(),
//     "age":   z.Int().GT(0).LT(150).Optional(),
// })
//
// // Define your struct
// type User struct {
//     Name  string `json:"name"`
//     Email string `json:"email"`
//     Age   *int   `json:"age,omitempty"`
// }
//
// // Use in route
// app.Post("/users",
//     func(ctx core.Context) error {
//         // Get parsed and validated body
//         user := usejsonbody.GetJsonBody(ctx).(*User)
//         // Or access directly from context
//         // user := ctx.Get("_jsonBody").(*User)
//
//         // Use validated user data
//         resp := core.NewResponse().
//             WithStatus(201).
//             WithBody(map[string]interface{}{
//                 "message": "User created",
//                 "user":    user,
//             })
//         return resp.Send(ctx)
//     },
//     usejsonbody.UseJsonBody(userSchema, &User{}),
// )
