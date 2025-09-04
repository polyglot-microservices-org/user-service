package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Contact represents the data model in MongoDB
type Contact struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name  string             `bson:"name" json:"name"`
    Phone string             `bson:"phone" json:"phone"`
}

var contactsCollection *mongo.Collection

// init connects to MongoDB
func init() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        mongoURI = "mongodb://user-db:27017"
    }

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    if err := client.Ping(ctx, nil); err != nil {
        log.Fatalf("Failed to ping MongoDB: %v", err)
    }

    fmt.Println("Connected to MongoDB successfully!")
    contactsCollection = client.Database("contacts_db").Collection("contacts")
}

// EnableCORS middleware
func EnableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// healthCheck handles the /healthz route for probes
func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// createContact handles POST /contacts
func createContact(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var contact Contact
    if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
        http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    if contact.Name == "" || contact.Phone == "" {
        http.Error(w, `{"error": "Missing name or phone"}`, http.StatusBadRequest)
        return
    }

    result, err := contactsCollection.InsertOne(context.TODO(), bson.M{
        "name":  contact.Name,
        "phone": contact.Phone,
    })
    if err != nil {
        http.Error(w, `{"error": "Failed to create contact"}`, http.StatusInternalServerError)
        return
    }

    contact.ID = result.InsertedID.(primitive.ObjectID)
    json.NewEncoder(w).Encode(bson.M{
        "message": "Contact created successfully",
        "contact": contact,
    })
}

// getContacts handles GET /contacts
func getContacts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var contacts []Contact
    cursor, err := contactsCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        http.Error(w, `{"error": "Failed to retrieve contacts"}`, http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var c Contact
        cursor.Decode(&c)
        contacts = append(contacts, c)
    }

    if err := cursor.Err(); err != nil {
        http.Error(w, `{"error": "Cursor error"}`, http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(contacts)
}

// getContact handles GET /contacts/{id}
func getContact(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id := r.URL.Path[len("/contacts/"):]
    if id == "" {
        http.Error(w, `{"error": "Missing contact ID"}`, http.StatusBadRequest)
        return
    }

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, `{"error": "Invalid contact ID"}`, http.StatusBadRequest)
        return
    }

    var c Contact
    err = contactsCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&c)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, `{"error": "Contact not found"}`, http.StatusNotFound)
            return
        }
        http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(c)
}

// updateContact handles PUT /contacts/{id}
func updateContact(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id := r.URL.Path[len("/contacts/"):]
    if id == "" {
        http.Error(w, `{"error": "Missing contact ID"}`, http.StatusBadRequest)
        return
    }

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, `{"error": "Invalid contact ID"}`, http.StatusBadRequest)
        return
    }

    var updateData map[string]string
    if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
        http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
        return
    }

    updateFields := bson.M{}
    if name, ok := updateData["name"]; ok {
        updateFields["name"] = name
    }
    if phone, ok := updateData["phone"]; ok {
        updateFields["phone"] = phone
    }

    result, err := contactsCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": updateFields})
    if err != nil {
        http.Error(w, `{"error": "Failed to update contact"}`, http.StatusInternalServerError)
        return
    }

    if result.MatchedCount == 0 {
        http.Error(w, `{"error": "Contact not found"}`, http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(bson.M{"message": "Contact updated successfully"})
}

// deleteContact handles DELETE /contacts/{id}
func deleteContact(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id := r.URL.Path[len("/contacts/"):]
    if id == "" {
        http.Error(w, `{"error": "Missing contact ID"}`, http.StatusBadRequest)
        return
    }

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, `{"error": "Invalid contact ID"}`, http.StatusBadRequest)
        return
    }

    result, err := contactsCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
    if err != nil {
        http.Error(w, `{"error": "Failed to delete contact"}`, http.StatusInternalServerError)
        return
    }

    if result.DeletedCount == 0 {
        http.Error(w, `{"error": "Contact not found"}`, http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(bson.M{"message": "Contact deleted successfully"})
}

func main() {
    router := http.NewServeMux()
    
    // Health check endpoint for Kubernetes probes
    router.HandleFunc("/healthz", healthCheck)

    // /contacts (no trailing slash)
    router.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "POST":
            createContact(w, r)
        case "GET":
            getContacts(w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    })

    // /contacts/{id}
    router.HandleFunc("/contacts/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/contacts/" && r.Method == "GET" {
            getContacts(w, r)
            return
        }

        switch r.Method {
        case "GET":
            getContact(w, r)
        case "PUT":
            updateContact(w, r)
        case "DELETE":
            deleteContact(w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    })

    handler := EnableCORS(router)

    port := os.Getenv("PORT")
    if port == "" {
        port = "5000"
    }

    fmt.Printf("Contacts API running on port %s...\n", port)
    log.Fatal(http.ListenAndServe(":"+port, handler))
}
