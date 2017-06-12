package main

import (
	"fmt"
	"log"
	"database/sql"
	"sync"
)

type userDeltas struct
{
	earnedMy float64
	earnedOther float64
	earnedRTB float64
	earnedSystem float64
	earnedReferral float64
	spent float64
}

type User struct {
	id uint32
	parentID uint32
	status uint16
	balance float64
	balanceLimit float64
	creditLimit float64
//	deltas userDeltas
	sync sync.Mutex
}

type UsersList map[uint32]*User
type Users struct {
	list UsersList
}

const USER_SQL = "SELECT id, parentID, status, balance, balanceLimit, creditLimit FROM user WHERE !(status&128)";

func (src *User)Reload(dst *User) {
	log.Println("Reload user ", src.id);
	src.parentID = dst.parentID;
	src.status = dst.status;
	src.balanceLimit = dst.balanceLimit;
	src.creditLimit = dst.creditLimit;
}

func (user *User)setFromRow(rows *sql.Rows) {
	err := rows.Scan(&user.id, &user.parentID, &user.status, &user.balance, &user.balanceLimit, &user.creditLimit)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *Users)LoadAll(db *sql.DB) {
	log.Print("Load users.")
	rows, err := db.Query(fmt.Sprintf("%s LIMIT 20000", USER_SQL))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	//log.Print(rows.affectedRows());
	r.list=make(UsersList)

	for rows.Next() {
		var user=new(User)
		user.setFromRow(rows)
		r.list[user.id] = user;
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(r.list), "users");
}

func (a *Users)Load(id uint32, db *sql.DB) *User {
	rows, err := db.Query(fmt.Sprintf("%s AND id=%d", USER_SQL, id))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		var user=new(User)
		user.setFromRow(rows)
		if a.list[user.id] != nil {
			a.list[user.id].Reload(user)
			//delete(ad)
		} else {
			log.Print("New user id=", user.id)
			a.list[user.id] = user
		}
		return a.list[user.id]
	} else {
		log.Println("Ad", id, "not found")
	}
	return nil;
}
