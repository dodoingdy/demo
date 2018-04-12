package pghandler

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"encoding/json"
	"strconv"
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

func UserInsert(username string) []byte{
	stmt, err := db.Prepare("INSERT INTO users(username) VALUES($1)")
	CheckErr(err)

	res, err := stmt.Exec(username)
	CheckErr(err)

	affect, err := res.RowsAffected()
	CheckErr(err)

	fmt.Println("rows affected:", affect)

	var result Users

	new_id, err := db.Query("select currval('uid_seq')")
	for new_id.Next() {
		var uid int
		err = new_id.Scan(&uid)
		CheckErr(err)
		result.Users = append(result.Users, User{Uid:uid, Username:username, Utype:"user"})
	}
	bytes, _ := json.Marshal(result)
	return bytes
}

func GetAllUsers() []byte {
	rows, err := db.Query("select * from users")
	CheckErr(err)

	var result Users

	for rows.Next() {
		var uid int
		var username string
		var utype string
		err = rows.Scan(&uid, &username, &utype)
		CheckErr(err)
		result.Users = append(result.Users, User{Uid:uid, Username:username, Utype:utype})
		fmt.Println("uid = ", uid, "username = ", username, "type = ", utype)
	}
	bytes, _ := json.Marshal(result)
	fmt.Println(bytes)
	return bytes
}

func GetALLRelationships(uid int) []byte {
	sql := "select userid,state,type from relationships where uid=" + strconv.Itoa(uid)
	fmt.Println(sql)
	rows, err := db.Query(sql)
	CheckErr(err)

	var result Relationships

	for rows.Next() {
		var userid int
		var state string
		var rtype string
		err := rows.Scan(&userid, &state, &rtype)
		CheckErr(err)
		result.Relationships = append(result.Relationships, Relationship{Userid:userid, State:state, Rtype:rtype})
		fmt.Println("userid = ", userid, "state = ", state, "type = ", rtype)
	}
	bytes, _ := json.Marshal(result)
	fmt.Println(string(bytes[:]))
	return bytes
}

func NewRelationship(uid int, user_id int, state string, rtype string) []byte{
	stmt, err := db.Prepare("INSERT INTO relationships(uid,userid,state,type) VALUES($1,$2,$3,$4)")
	CheckErr(err)

	sql := "select state from relationships where uid=" + strconv.Itoa(user_id) + " and userid=" + strconv.Itoa(uid)

	var result Relationships

	if state == "liked" {
		row, err := db.Query(sql)
		CheckErr(err)
		fmt.Println(row)

		var oldstate string
		key := 0

		for row.Next() {
			err := row.Scan(&oldstate)
			CheckErr(err)

			if oldstate == "liked" {
				key = 1
				upstmt, err := db.Prepare("update relationships set state=$1 where uid=$2 and userid=$3")
				CheckErr(err)

				res, err := upstmt.Exec("matched", user_id, uid)
				CheckErr(err)

				affect, err := res.RowsAffected()
				CheckErr(err)
									
				fmt.Println("rows affected:", affect)

				res, err = stmt.Exec(uid, user_id, "matched", rtype)
				CheckErr(err)

				affect, err = res.RowsAffected()
				CheckErr(err)
	
				fmt.Println("rows affected:", affect)
				result.Relationships = append(result.Relationships, Relationship{Userid:user_id, State:"matched", Rtype:rtype})
			}
		}
		if key == 0 {
			res, err := stmt.Exec(uid, user_id, state, rtype)
			CheckErr(err)

			affect, err := res.RowsAffected()
			CheckErr(err)
	
			fmt.Println("rows affected:", affect)
			result.Relationships = append(result.Relationships, Relationship{Userid:user_id, State:state, Rtype:rtype})
		}
	} else {
		res, err := stmt.Exec(uid, user_id, state, rtype)
		CheckErr(err)

		affect, err := res.RowsAffected()
		CheckErr(err)
	
		fmt.Println("rows affected:", affect)
		result.Relationships = append(result.Relationships, Relationship{Userid:user_id, State:state, Rtype:rtype})
	}
	bytes, _ := json.Marshal(result)
	fmt.Println(string(bytes[:]))
	return bytes
}
