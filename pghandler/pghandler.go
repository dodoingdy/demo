package pghandler

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"encoding/json"
	"strconv"
	"log"
)

var db *sql.DB

type Users struct{
	Users []User `json:users`
}

type User struct {
	Uid int `json:uid`
	Username string `json:username`
	Utype string `json:utype`
}

type Relationships struct{
	Relationships []Relationship `json:relationships`
}

type Relationship struct {
	Userid int `json:userid`
	State string `json:state`
	Rtype string `json:rtype`
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func SqlOpen() {
	var err error
	db, err = sql.Open("postgres", "user=dingyi password=dy20130429 dbname=test sslmode=disable")
	CheckErr(err)
}

func SqlClose(){
	db.Close()
}

func(result Users) UserInsert(username string) []byte{
	stmt, err := db.Prepare("INSERT INTO users(username) VALUES($1)")
	CheckErr(err)

	res, err := stmt.Exec(username)
	CheckErr(err)

	affect, err := res.RowsAffected()
	CheckErr(err)

	log.Println("rows affected:", affect)

	newId, err := db.Query("select currval('uid_seq')")
	for newId.Next() {
		var uid int
		err = newId.Scan(&uid)
		CheckErr(err)
		result.Users = append(result.Users, User{Uid:uid, Username:username, Utype:"user"})
	}
	bytes, _ := json.Marshal(result)
	return bytes
}

func(result Users) GetAllUsers() []byte {
	rows, err := db.Query("select * from users")
	CheckErr(err)

	for rows.Next() {
		var uid int
		var username string
		var utype string
		err = rows.Scan(&uid, &username, &utype)
		CheckErr(err)
		result.Users = append(result.Users, User{Uid:uid, Username:username, Utype:utype})
		log.Println("uid = ", uid, "username = ", username, "type = ", utype)
	}
	bytes, _ := json.Marshal(result)
	return bytes
}

func(result Relationships) GetALLRelationships(uid int) []byte {
	sql := "select userid,state,type from relationships where uid=" + strconv.Itoa(uid)
	log.Println(sql)
	rows, err := db.Query(sql)
	CheckErr(err)

	for rows.Next() {
		var userid int
		var state string
		var rtype string
		err := rows.Scan(&userid, &state, &rtype)
		CheckErr(err)
		result.Relationships = append(result.Relationships, Relationship{Userid:userid, State:state, Rtype:rtype})
		log.Println("userid = ", userid, "state = ", state, "type = ", rtype)
	}
	bytes, _ := json.Marshal(result)
	//fmt.Println(string(bytes[:]))
	return bytes
}

func(result Relationships) NewRelationship(uid int, userId int, state string, rtype string) []byte{
	stmt, err := db.Prepare("INSERT INTO relationships(uid,userid,state,type) VALUES($1,$2,$3,$4)")
	CheckErr(err)

	sql := "select state from relationships where uid=" + strconv.Itoa(userId) + " and userid=" + strconv.Itoa(uid)

	if state == "liked" {
		row, err := db.Query(sql)
		CheckErr(err)
		log.Println(row)

		var oldstate string
		key := 0

		for row.Next() {
			err := row.Scan(&oldstate)
			CheckErr(err)

			if oldstate == "liked" {
				key = 1
				upstmt, err := db.Prepare("update relationships set state=$1 where uid=$2 and userid=$3")
				CheckErr(err)

				res, err := upstmt.Exec("matched", userId, uid)
				CheckErr(err)

				affect, err := res.RowsAffected()
				CheckErr(err)
									
				fmt.Println("rows affected:", affect)

				res, err = stmt.Exec(uid, userId, "matched", rtype)
				CheckErr(err)

				affect, err = res.RowsAffected()
				CheckErr(err)
	
				log.Println("rows affected:", affect)
				result.Relationships = append(result.Relationships, Relationship{Userid:userId, State:"matched", Rtype:rtype})
			}
		}
		if key == 0 {
			res, err := stmt.Exec(uid, userId, state, rtype)
			CheckErr(err)

			affect, err := res.RowsAffected()
			CheckErr(err)
	
			log.Println("rows affected:", affect)
			result.Relationships = append(result.Relationships, Relationship{Userid:userId, State:state, Rtype:rtype})
		}
	} else {
		var oldstate string
		err := db.QueryRow(sql).Scan(&oldstate)
		CheckErr(err)

		if oldstate == "matched"{
			upstmt, err := db.Prepare("update relationships set state=$1 where uid=$2 and userid=$3")
			CheckErr(err)

			res, err := upstmt.Exec("liked", userId, uid)
			CheckErr(err)

			affect, err := res.RowsAffected()
			CheckErr(err)
									
			fmt.Println("rows affected:", affect)

			res, err = upstmt.Exec(state, uid, userId)
			CheckErr(err)

			affect, err = res.RowsAffected()
			CheckErr(err)
	
			log.Println("rows affected:", affect)
			result.Relationships = append(result.Relationships, Relationship{Userid:userId, State:state, Rtype:rtype})
		} else {
			res, err := stmt.Exec(uid, userId, state, rtype)
			CheckErr(err)

			affect, err := res.RowsAffected()
			CheckErr(err)
	
			log.Println("rows affected:", affect)
			result.Relationships = append(result.Relationships, Relationship{Userid:userId, State:state, Rtype:rtype})
		}
	}
	bytes, _ := json.Marshal(result)
	log.Println(string(bytes[:]))
	return bytes
}
