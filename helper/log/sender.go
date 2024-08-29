package log

import (
	"api/helper/atdb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func LogSenderReceiverUpdate(sender, receiver string, db *mongo.Database) {
	const logcollection = "logsent"
	filter := primitive.M{
		"sender":   "",
		"receiver": "",
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
