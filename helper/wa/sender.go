package wa

import (
	"api/helper/atdb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func LogSenderReceiverUpdate(sender, receiver string, db *mongo.Database) {
	const logcollection = "logsent"
	filter := primitive.M{
		"sender":   sender,
		"receiver": receiver,
	}
	_, err := atdb.GetOneDoc[LogSenderReceiver](db, logcollection, filter)
	if err == mongo.ErrNoDocuments {
		newdoc := LogSenderReceiver{
			Sender:   sender,
			Receiver: receiver,
		}
		atdb.InsertOneDoc(db, logcollection, newdoc)
	}
}

func GetSenderNumber(receiver string, db *mongo.Database) string {
	doc, err := atdb.GetOneDoc[LogSenderReceiver](db, "logsent", primitive.M{"receiver": receiver})
	if err == mongo.ErrNoDocuments {
		return ""
	}
	return doc.Sender
}

func GetOfficialSenderNumber(user string, db *mongo.Database) string {
	doc, _ := atdb.GetOneLowestDoc[LogSenderCounterUsage](db, "senderofficial", primitive.M{"users": user}, "counter")
	doc.Counter += 1
	atdb.ReplaceOneDoc(db, "senderofficial", primitive.M{"_id": doc.ID}, doc)
	return doc.Sender
}
