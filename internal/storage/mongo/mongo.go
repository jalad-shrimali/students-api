package mongo

import (
        "context"
        "fmt"
        "log"
        "os"
        "strconv"
        "time"

        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/bson/primitive"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"

        "github.com/jalad-shrimali/students-api/internal/types"
)

type Mongo struct {
        Client         *mongo.Client
        DbName         string
        CollectionName string
}

func New(uri string, dbName string, collectionName string) (*Mongo, error) {
        clientOptions := options.Client().ApplyURI(uri)
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        client, err := mongo.Connect(ctx, clientOptions)
        if err != nil {
                return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
        }

        err = client.Ping(ctx, nil)
        if err != nil {
                return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
        }

        log.Println("Connected to MongoDB!")

        return &Mongo{
                Client:         client,
                DbName:         dbName,
                CollectionName: collectionName,
        }, nil
}

func NewFromEnv() (*Mongo, error) {
        uri := os.Getenv("MONGO_URI")
        dbName := os.Getenv("MONGO_DB")
        collectionName := os.Getenv("MONGO_COLLECTION")

        if uri == "" || dbName == "" || collectionName == "" {
                return nil, fmt.Errorf("MONGO_URI, MONGO_DB, or MONGO_COLLECTION environment variables not set")
        }

        return New(uri, dbName, collectionName)
}

func (m *Mongo) getCollection() *mongo.Collection {
        return m.Client.Database(m.DbName).Collection(m.CollectionName)
}

func (m *Mongo) CreateStudent(name string, age int, email string) (int64, error) {
        student := types.Student{
                Name:  name,
                Age:   age,
                Email: email,
        }

        collection := m.getCollection()
        result, err := collection.InsertOne(context.Background(), student)
        if err != nil {
                return 0, fmt.Errorf("failed to insert student: %w", err)
        }
        oid, ok := result.InsertedID.(primitive.ObjectID)
        if !ok {
                return 0, fmt.Errorf("failed to convert object ID to primitive.ObjectID")
        }
        //Convert to int64, while might lose some data, it matches the interface.
        id, err := strconv.ParseInt(oid.Hex()[18:],16,64)
        if err != nil {
                return 0, fmt.Errorf("failed to convert object ID to int64: %w", err)
        }

        return id, nil
}

func (m *Mongo) GetStudentById(id int64) (types.Student, error) {
        collection := m.getCollection()
        oid, err := primitive.ObjectIDFromHex(fmt.Sprintf("%018x%016x", 0, id))
        if err != nil {
                return types.Student{}, fmt.Errorf("invalid object id: %w", err)
        }
        filter := bson.M{"_id": oid}
        var student types.Student
        err = collection.FindOne(context.Background(), filter).Decode(&student)
        if err != nil {
                if err == mongo.ErrNoDocuments {
                        return types.Student{}, fmt.Errorf("no student found with id %d", id)
                }
                return types.Student{}, fmt.Errorf("failed to find student: %w", err)
        }
        id, err = strconv.ParseInt(oid.Hex()[18:], 16, 64)
        if err != nil {
                return types.Student{}, fmt.Errorf("failed to convert object ID to int64: %w", err)
        }
        if err != nil {
                return types.Student{}, fmt.Errorf("failed to convert object ID to int64: %w", err)
        }
        student.Id = id
        return student, nil
}

func (m *Mongo) GetAllStudents() ([]types.Student, error) {
        collection := m.getCollection()
        cursor, err := collection.Find(context.Background(), bson.M{})
        if err != nil {
                return nil, fmt.Errorf("failed to find students: %w", err)
        }
        defer cursor.Close(context.Background())

        var students []types.Student
        for cursor.Next(context.Background()) {
                var student types.Student
                if err := cursor.Decode(&student); err != nil {
                        return nil, fmt.Errorf("failed to decode student: %w", err)
                }
                oid := cursor.Current.Lookup("_id").ObjectID()
                id, err := strconv.ParseInt(oid.Hex()[18:], 16, 64)
                if err != nil {
                        return nil, fmt.Errorf("failed to convert object ID to int64: %w", err)
                }
                student.Id = id
                students = append(students, student)
        }
        if err := cursor.Err(); err != nil {
                return nil, fmt.Errorf("cursor error: %w", err)
        }
        return students, nil
}

func (m *Mongo) UpdateStudent(id int64, name string, age int, email string) (types.Student, error) {
        collection := m.getCollection()
        oid, err := primitive.ObjectIDFromHex(fmt.Sprintf("%018x%016x", 0, id))
        if err != nil {
                return types.Student{}, fmt.Errorf("invalid object id: %w", err)
        }
        filter := bson.M{"_id": oid}
        update := bson.M{"$set": bson.M{"name": name, "age": age, "email": email}}
        _, err = collection.UpdateOne(context.Background(), filter, update)
        if err != nil {
                return types.Student{}, fmt.Errorf("failed to update student: %w", err)
        }
        return types.Student{
                Id:    id,
                Name:  name,
                Age:   age,
                Email: email,
        }, nil
}

func (m *Mongo) DeleteStudent(id int64) error {
        collection := m.getCollection()
        oid, err := primitive.ObjectIDFromHex(fmt.Sprintf("%018x%016x", 0, id))
        if err != nil {
                return fmt.Errorf("invalid object id: %w", err)
        }
        filter := bson.M{"_id": oid}
        _, err = collection.DeleteOne(context.Background(), filter)
        if err != nil {
                return fmt.Errorf("failed to delete student: %w", err)
        }
        return nil
}