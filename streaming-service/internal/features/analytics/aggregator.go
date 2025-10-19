package analytics

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetDailyAnalyticsPipeline() mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{{
			Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "date", Value: bson.D{
						{Key: "$dateToString", Value: bson.D{
							{Key: "format", Value: "%Y-%m-%d"},
							{Key: "date", Value: "$created_at"},
							{Key: "timezone", Value: "UTC"},
						}},
					}},
					{Key: "topic", Value: "$topic"},
				}},
				{Key: "session_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "all_participants", Value: bson.D{{Key: "$push", Value: "$participants.user_id"}}},
			},
		}},

		bson.D{{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "date", Value: "$_id.date"},
				{Key: "topic", Value: "$_id.topic"},
				{Key: "session_count", Value: 1},
				{Key: "participants_set", Value: bson.D{
					{Key: "$setUnion", Value: bson.D{
						{Key: "$reduce", Value: bson.D{
							{Key: "input", Value: "$all_participants"},
							{Key: "initialValue", Value: bson.A{}},
							{Key: "in", Value: bson.D{
								{Key: "$concatArrays", Value: bson.A{"$$value", "$$this"}},
							}},
						}},
					}},
				}},
				{Key: "participant_count", Value: bson.D{
					{Key: "$size", Value: bson.D{
						{Key: "$setUnion", Value: bson.D{
							{Key: "$reduce", Value: bson.D{
								{Key: "input", Value: "$all_participants"},
								{Key: "initialValue", Value: bson.A{}},
								{Key: "in", Value: bson.D{
									{Key: "$concatArrays", Value: bson.A{"$$value", "$$this"}},
								}},
							}},
						}},
					}},
				}},
			},
		}},

		bson.D{{
			Key: "$sort", Value: bson.D{
				{Key: "date", Value: 1},
				{Key: "participant_count", Value: -1},
			},
		}},

		bson.D{{
			Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$date"},
				{Key: "session_count", Value: bson.D{{Key: "$sum", Value: "$session_count"}}},
				{Key: "top_topic", Value: bson.D{{Key: "$first", Value: "$topic"}}},
				{Key: "top_topic_participant_count", Value: bson.D{{Key: "$first", Value: "$participant_count"}}},
				{Key: "all_topic_participants", Value: bson.D{{Key: "$push", Value: "$participants_set"}}},
			},
		}},

		bson.D{{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "created_at", Value: "$_id"},
				{Key: "session_count", Value: 1},
				{Key: "top_topic", Value: 1},
				{Key: "top_topic_participant_count", Value: 1},
				{Key: "participant_count", Value: bson.D{
					{Key: "$size", Value: bson.D{
						{Key: "$setUnion", Value: bson.D{
							{Key: "$reduce", Value: bson.D{
								{Key: "input", Value: "$all_topic_participants"},
								{Key: "initialValue", Value: bson.A{}},
								{Key: "in", Value: bson.D{
									{Key: "$concatArrays", Value: bson.A{"$$value", "$$this"}},
								}},
							}},
						}},
					}},
				}},
			},
		}},

		bson.D{{
			Key: "$sort", Value: bson.D{
				{Key: "created_at", Value: 1},
			},
		}},
	}
}
