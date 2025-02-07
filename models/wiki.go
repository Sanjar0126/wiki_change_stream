package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WikiRecentChangesMeta struct {
	URI       string    `json:"uri" bson:"uri"`
	RequestID string    `json:"request_id" bson:"request_id"`
	ID        string    `json:"id" bson:"id"`
	Dt        time.Time `json:"dt" bson:"dt"`
	Domain    string    `json:"domain" bson:"domain"`
	Stream    string    `json:"stream" bson:"stream"`
	Topic     string    `json:"topic" bson:"topic"`
	Partition int       `json:"partition" bson:"partition"`
	Offset    int       `json:"offset" bson:"offset"`
}

type WikiRecentChangesLength struct {
	New int `json:"new" bson:"new"`
}

type Revision struct {
	New int64 `json:"new" bson:"new"`
}

type WikiRecentChanges struct {
	BId              primitive.ObjectID      `json:"_id" bson:"_id"`         //nolint
	Schema           string                  `json:"$schema" bson:"$schema"` //nolint
	Meta             WikiRecentChangesMeta   `json:"meta" bson:"meta"`
	ID               int64                   `json:"id" bson:"id"`
	Type             string                  `json:"type" bson:"type"`
	Namespace        int                     `json:"namespace" bson:"namespace"`
	Title            string                  `json:"title" bson:"title"`
	TitleURL         string                  `json:"title_url" bson:"title_url"`
	Comment          string                  `json:"comment" bson:"comment"`
	Timestamp        int                     `json:"timestamp" bson:"timestamp"`
	User             string                  `json:"user" bson:"user"`
	Bot              bool                    `json:"bot" bson:"bot"`
	NotifyURL        string                  `json:"notify_url" bson:"notify_url"`
	Minor            bool                    `json:"minor" bson:"minor"`
	Patrolled        bool                    `json:"patrolled" bson:"patrolled"`
	Length           WikiRecentChangesLength `json:"length" bson:"length"`
	Revision         Revision                `json:"revision" bson:"revision"`
	ServerURL        string                  `json:"server_url" bson:"server_url"`
	ServerName       string                  `json:"server_name" bson:"server_name"`
	ServerPrefix     string                  `json:"server_prefix" bson:"server_prefix"`
	ServerScriptPath string                  `json:"server_script_path" bson:"server_script_path"`
	Wiki             string                  `json:"wiki" bson:"wiki"`
	Parsedcomment    string                  `json:"parsedcomment" bson:"parsedcomment"`
}
