package handlers

import (
	"context"

	"polygon-server/models"
	"polygon-server/response"
	"polygon-server/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB *mongo.Database
}

// CheckPassword checks if the provided password matches the hashed password
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	var users []models.User

	// Define the projection to exclude the password field
	projection := bson.D{
		{"password", 0},
	}

	// Use the projection in the FindOptions
	findOptions := options.Find().SetProjection(projection)

	cursor, err := h.DB.Collection("users").Find(context.Background(), bson.D{}, findOptions)
	if err != nil {
		return response.InternalServerError(c, err, "Error finding users")
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &users); err != nil {
		return response.InternalServerError(c, err, "Error decoding users")
	}

	return response.SuccessHandler(c, fiber.StatusOK, "Users retrieved successfully", users)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return response.BadRequest(c, err, "Cannot parse JSON")
	}

	isValid, validationErrors := utils.ValidateStruct(user)
	if !isValid {
		return response.ValidationError(c, validationErrors)
	}

	// Check if the user already exists based on email
	var existingUser models.User
	err := h.DB.Collection("users").FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return response.BadRequest(c, nil, "User with this email already exists")
	}

	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return response.InternalServerError(c, err, "Error hashing password")
	}
	user.Password = hashedPassword

	user.ID = primitive.NewObjectID()
	_, err = h.DB.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return response.InternalServerError(c, err, "Error inserting user")
	}

	// Remove password from response
	user.Password = ""
	return response.SuccessHandler(c, fiber.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	input := new(models.User)
	if err := c.BodyParser(input); err != nil {
		return response.BadRequest(c, err, "Cannot parse JSON")
	}

	var user models.User
	err := h.DB.Collection("users").FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return response.InternalServerError(c, err, "Invalid email or password")
	}

	if !CheckPassword(user.Password, input.Password) {
		return response.InternalServerError(c, nil, "Invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		return response.InternalServerError(c, err, "Failed to generate token")
	}

	return response.SuccessHandler(c, fiber.StatusOK, "Login successful", fiber.Map{"token": token})
}
