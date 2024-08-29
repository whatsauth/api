package log

import "go.mongodb.org/mongo-driver/bson/primitive"

type LogSenderReceiver struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Sender   string             `bson:"sender"`
	Receiver string             `bson:"receiver"`
}

type LogSenderCounterUsage struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Sender  string             `bson:"sender"`
	Counter uint64             `bson:"counter,omitempty"`
	Users   []string           `bson:"users,omitempty"`
}
