package mongohandler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// InitialDatabase - InitialDatabase
func (c *Mgo) InitialDatabase() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:   []string{MongoHost},
		Timeout: 60 * time.Second,
		//PoolLimit: 100000,
	}
	_db, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		fmt.Println("MONGO ERROR: ", err)
	}
	db = _db
}

// Close - Close connection to mongoDB
func (c *Mgo) Close() {
	db.Close()
}

// CountDocuments - CountDocuments
func (c *Mgo) CountDocuments(MongoHost, DBName, collection string) (int, error) {
	n, err := db.DB(DBName).C(collection).Count()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	return n, nil
}

// GetCollectionNames - GetCollectionNames
func (c *Mgo) GetCollectionNames() ([]string, error) {
	n, err := db.DB(DBName).CollectionNames()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return n, nil
}

//FindByID - Find one document in Mongo DB by ID
func (c *Mgo) FindByID(MongoHost, DBName, collection string, id string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := db.DB(DBName).C(collection).Find(bson.M{"_id": id}).One(&result)

	if err == mgo.ErrNotFound || err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

//FindAll - Find all document in Mongo DB
func (c *Mgo) FindAll(MongoHost, DBName, collection string) []map[string]interface{} {
	var result []map[string]interface{}
	err := db.DB(DBName).C(collection).Find(bson.M{}).All(&result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}

//FindAllID - Find all Document's ID in Mongo DB
func (c *Mgo) FindAllID(MongoHost, DBName, collection string) []map[string]interface{} {
	var result []map[string]interface{}
	err := db.DB(DBName).C(collection).Find(bson.M{}).Select(bson.M{"_id": "1"}).All(&result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}

//FindOneMongoByField - Find one document in Mongo DB by any Field
func (c *Mgo) FindOneMongoByField(MongoHost, DBName, collection, f, v string) (map[string]interface{}, bool) {
	var result map[string]interface{}
	selector := bson.M{f: bson.M{"$regex": v}}
	err := db.DB(DBName).C(collection).Find(selector).One(result)

	if err == mgo.ErrNotFound {
		return nil, false
	}
	return result, true
}

//SaveMongo - Save Data to Mongo DB
func (c *Mgo) SaveMongo(MongoHost, DBName, collection string, ID string, data map[string]interface{}) {
	_, err := db.DB(DBName).C(collection).Upsert(bson.M{"_id": ID}, data)
	if err != nil {
		fmt.Println(err)
	}
}

//UpdateMongo - UpdateMongo
func (c *Mgo) UpdateMongo(MongoHost, DBName, collection, id string, field, value string) error {
	// value := make(map[string]interface{})
	// err := json.Unmarshal(data, &value)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	filter := bson.M{"_id": bson.M{"$eq": id}}
	update := bson.M{"$set": bson.M{field: value}}

	err := db.DB(DBName).C(collection).Update(filter, update)
	return err
}

//GetDataInOrderedByField - GetDataInOrderedByField
func (c *Mgo) GetDataInOrderedByField(MongoHost, DBName, collection, fields string) (interface{}, bool) {
	var result []map[string]interface{}

	err := db.DB(DBName).C(collection).Find(nil).Sort("-" + fields).All(&result)
	if err != nil || len(result) == 0 {
		fmt.Println(err)
		return nil, false
	}

	return result, true
}

func isNumDot(s string) bool {
	dotFound := false
	for _, v := range s {
		if v == '.' {
			if dotFound {
				return false
			}
			dotFound = true
		} else if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

// SearchInMongoByRange - Search Data By Range in MongoDB
func (c *Mgo) SearchInMongoByRange(MongoHost, DBName, bucket string, bodyBytes []byte) (interface{}, bool) {
	var bodyMap map[string]string
	json.Unmarshal(bodyBytes, &bodyMap)

	if bodyMap["field"] == "" {
		fmt.Println("ERR: Field didn't Fill")
		return nil, false
	}
	field := bodyMap["field"]

	if bodyMap["gte"] == "" && bodyMap["gt"] == "" && bodyMap["lte"] == "" && bodyMap["lt"] == "" {
		fmt.Println("ERR: GTE, LTE, GT, LT didn't Fill")
		return nil, false
	}

	if bodyMap["gte"] != "" && bodyMap["gt"] != "" {
		return nil, false
	}
	if bodyMap["lte"] != "" && bodyMap["lt"] != "" {
		return nil, false
	}

	var andQuery []bson.M
	for key, val := range bodyMap {
		if key == "gte" || key == "gt" || key == "lte" || key == "lt" {

			var selector bson.M
			if isNumDot(val) != true {
				selector = bson.M{field: bson.M{"$" + key: val}}
			} else {
				valInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					fmt.Println(err)
					return nil, false
				}
				selector = bson.M{field: bson.M{"$" + key: valInt}}

			}
			andQuery = append(andQuery, selector)
		}
	}

	var retValue []bson.M
	err := db.DB(DBName).C(bucket).Find(bson.M{"$and": andQuery}).All(&retValue)
	if err != nil || len(retValue) == 0 {
		fmt.Println(err)
		return nil, false
	}

	return retValue, true
}

//DeleteInMongo - Delete Data in MongoDB
func (c *Mgo) DeleteInMongo(MongoHost, DBName, bucket, deleteByID string) error {
	err := db.DB(DBName).C(bucket).Remove(bson.M{"_id": deleteByID})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//DropCollectionInMongo - Delete a Collection in MongoDB
func (c *Mgo) DropCollectionInMongo(MongoHost, DBName, bucket string) {
	err := db.DB(DBName).C(bucket).DropCollection()
	if err != nil {
		fmt.Println(err)
	}
}

//DropManyCollectionInMongo - DropManyCollectionInMongo
func (c *Mgo) DropManyCollectionInMongo(MongoHost, DBName string, listCollection []string) {
	for _, collection := range listCollection {
		collection = strings.ReplaceAll(collection, "$", "@")
		err := db.DB(DBName).C(collection).DropCollection()
		if err != nil {
			fmt.Println("Error while drop collection ", collection)
		}
	}
}

// DropDatabase - Drop database
func (c *Mgo) DropDatabase(MongoHost, DBName string) {
	err := db.DB(DBName).DropDatabase()
	if err != nil {
		fmt.Println("Error while drop database ", DBName)
	}
}

// FindAllInMongo - FindAllInMongo
func (c *Mgo) FindAllInMongo(collection string) []string {
	var result []bson.M
	var retData []string
	err := db.DB(DBName).C(collection).Find(bson.M{}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	for _, r := range result {
		rTmp := fmt.Sprintf("%s", r["_id"])
		retData = append(retData, rTmp)
	}
	return retData
}

//SortInMongo - SortInMongo
func (c *Mgo) SortInMongo(collection, sortByField string) (interface{}, error) {
	var result []bson.M
	err := db.DB(DBName).C(collection).Find(bson.M{}).Sort(sortByField).All(&result)
	if err != nil || len(result) == 0 {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

//LimitSortInMongo - LimitSortInMongo
func (c *Mgo) LimitSortInMongo(collection, sortByField string, limit int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	if sortByField == "" {
		err := db.DB(DBName).C(collection).Find(bson.M{}).Limit(limit).All(&result)
		if err != nil || len(result) == 0 {
			fmt.Println(err)
			return nil, err
		}
	} else {
		err := db.DB(DBName).C(collection).Find(bson.M{}).Sort(sortByField).All(&result)
		if err != nil || len(result) == 0 {
			fmt.Println(err)
			return nil, err
		}
	}
	return result, nil
}

//PaginateWithSkip - Phân trang dữ liệu.
// VD: Page 1, Limit 2 => Page 1: record 1, record 2; Page 2: record 3, record 4
//	   Page 2, Limit 3 => Page 1: record 1, record 2, record 3; Page 2: record 4, record 5, record 6
func (c *Mgo) PaginateWithSkip(collection string, page, limit int) ([]interface{}, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("Limit must > 0")
	}

	start := time.Now()
	skip := (page - 1) * limit

	records := make([]interface{}, Limit)
	err := db.DB(DBName).C(collection).Find(nil).Sort("_id").Skip(skip).Limit(limit).All(&records)
	if err != nil || len(records) == 0 {
		fmt.Println(err)
		return nil, err
	}
	rFirstMap := records[0].(bson.M)
	rLastMap := records[len(records)-1].(bson.M)
	rFirst := fmt.Sprintf("%s", rFirstMap["_id"])
	rLast := fmt.Sprintf("%s", rLastMap["_id"])

	fmt.Printf("paginateWithSkip -> page %d with record from %s to %s in %s\n", page, rFirst, rLast, time.Since(start))
	return records, nil
}

//FindAllRegexByID - Find All Data from Mongo DB by ID and Regex
func (c *Mgo) FindAllRegexByID(MongoHost, DBName, collection, id string) []map[string]interface{} {
	var result []map[string]interface{}
	err := db.DB(DBName).C(collection).Find(bson.M{"_id": bson.M{"$regex": id}}).Sort("timestamp").All(&result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}