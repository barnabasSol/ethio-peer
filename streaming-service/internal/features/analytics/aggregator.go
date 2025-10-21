package analytics

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetHourlyAnalyticsPipeline() mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "participant_count", Value: bson.D{
					{Key: "$size", Value: "$participants"},
				}},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "$dateToString", Value: bson.D{
						{Key: "format", Value: "%Y-%m-%d %H:00"},
						{Key: "date", Value: "$created_at"},
						{Key: "timezone", Value: "Africa/Addis_Ababa"},
					}},
				}},
				{Key: "total_participants", Value: bson.D{
					{Key: "$sum", Value: "$participant_count"},
				}},
				{Key: "sessions_created", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "hour", Value: "$_id"},
				{Key: "sessions_created", Value: 1},
				{Key: "total_participants", Value: 1},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "hour", Value: 1},
			}},
		},
	}
}

func GetDailyAnalyticsPipeline() mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "participant_count", Value: bson.D{
					{Key: "$size", Value: "$participants"},
				}},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "$dateToString", Value: bson.D{
						{Key: "format", Value: "%Y-%m-%d"},
						{Key: "date", Value: "$created_at"},
					}},
				}},
				{Key: "total_participants", Value: bson.D{
					{Key: "$sum", Value: "$participant_count"},
				}},
				{Key: "sessions_created", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "date", Value: "$_id"},
				{Key: "sessions_created", Value: 1},
				{Key: "total_participants", Value: 1},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "date", Value: 1},
			}},
		},
	}
}
