package proto

// 此文件为gRPC proto文件的Go代码生成版本
// 实际使用时需要使用protoc命令生成
// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto

// ProductResponse 商品响应
type ProductResponse struct {
	Id          uint64  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
	Category    string  `json:"category"`
	ImageUrl    string  `json:"image_url"`
	Status      int32   `json:"status"`
}

// UserResponse 用户响应
type UserResponse struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
