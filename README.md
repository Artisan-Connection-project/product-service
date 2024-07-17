# product-service
# need to edite mannually 
<!-- type Product struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" db:"id,omitempty"`
	Name        string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" db:"name,omitempty"`
	Description string  `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty" db:"description,omitempty"`
	Price       float64 `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty" db:"price,omitempty"`
	CategoryId  string  `protobuf:"bytes,5,opt,name=category_id,json=categoryId,proto3" json:"category_id,omitempty" db:"category_id,omitempty"`
	ArtisanId   string  `protobuf:"bytes,6,opt,name=artisan_id,json=artisanId,proto3" json:"artisan_id,omitempty" db:"artisan_id"`
	Quantity    int32   `protobuf:"varint,7,opt,name=quantity,proto3" json:"quantity,omitempty" db:"quantity"`
	CreatedAt   string  `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   string  `protobuf:"bytes,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" db:"updated_at"`
} -->

<!-- type Rating struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" db:"id,omitempty"`
	ProductId string  `protobuf:"bytes,2,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty" db:"product_id,omitempty"`
	UserId    string  `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" db:"user_id,omitempty"`
	Rating    float64 `protobuf:"fixed64,4,opt,name=rating,proto3" json:"rating,omitempty" db:"rating,omitempty"`
	Comment   string  `protobuf:"bytes,5,opt,name=comment,proto3" json:"comment,omitempty" db:"comment,omitempty"`
	CreatedAt string  `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty" db:"created_at,omitpty"`
} -->
<!-- type Order struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" db:"id,omitempty"`
	UserId          string  `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" db:"user_id,omitempty"`
	TotalAmount     float64 `protobuf:"fixed64,3,opt,name=total_amount,json=totalAmount,proto3" json:"total_amount,omitempty" db:"total_amount"`
	Status          string  `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty" db:"status"`
	ShippingAddress string  `protobuf:"bytes,5,opt,name=shipping_address,json=shippingAddress,proto3" json:"shipping_address,omitempty" db:"shippingAddres"`
	CreatedAt       string  `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty" db:"created_at"`
	UpdatedAt       string  `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" db:"updated_at"`
} -->
