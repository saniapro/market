package main

import (
	"fmt"
	"log"
	"database/sql"
	"sync"
)

//type extType

type TargetDeltas struct 
{
	views uint32
	clicks uint32
	spent float64
}


type Target struct {
	id uint32
	status uint16
	sync sync.Mutex
	priority uint8
	aclHours uint32
	aclWeekDay uint8
	aclRegion uint32
	aclAgeSex uint8
	aclCats uint64
	retargetID uint32
	viewsUserDayLimit uint8
	viewsDayLimit  uint32
	viewsToday uint32
	viewsYesterday uint32
	views uint32
	clicksToday uint32
	clicksYesterday uint32
	clicks uint32
	spentLimit float64
	spent float64
	spentDayLimit float64
	spentToday float64
	maxPrice float64

	deltas TargetDeltas
	_5minViews intc
	_5minViewsLimit intc

	_5minSpent float64
	_5minSpentLimit float64
}

type TargetsList map[uint32]*Target

type Targets struct {
	list TargetsList
}

const TARGET_SQL_FIELDS = "t.status AS tstatus, priority, aclHours, aclWeekDay, aclRegion, aclAgeSex, aclCats, retargetID, viewsUserDayLimit, periods, viewsDayLimit, viewsToday, viewsYesterday, views, clicksToday, clicksYesterday, clicks, spentLimit, spent, spentDayLimit, spentToday, maxPrice"

/*func (ad *Ad)setFromRow(rows *sql.Rows) {
	err := rows.Scan(&ad.id, &ad.campaign, &ad.target, &ad.status, &ad.width, &ad.height, &ad.ext, &ad.text, &ad.text1, &ad.price)
	if err != nil {
		log.Fatal(err)
	}
}

func (src *Ad)Reload(dst *Ad) {
	log.Println("Reload ", src.id);
	src.campaign = dst.campaign;
	src.status = dst.status;
	src.width = dst.width;
	src.height = dst.height;
	src.ext = dst.ext;
	src.price = dst.price;

	src.sync.Lock()
	defer src.sync.Unlock()

	src.text = dst.text;
	src.text1 = dst.text1;

//	_status = src->_status;
//	_format = src->_format;
//	_catID = src->_catID;
//	_picSizeStatus = src->_picSizeStatus;
//	_version = src->_version;
}

func (a *Ads)Load(id uint32, db *sql.DB) *Ad {
	rows, err := db.Query(fmt.Sprintf("%s AND id=%d", AD_SQL, id))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		var ad=new(Ad)
		ad.setFromRow(rows)
		if a.list[ad.id] != nil {
			a.list[ad.id].Reload(ad)
			//delete(ad)
		} else {
			log.Print("New ad id=", ad.id)
			a.list[ad.id] = ad
		}
		return a.list[ad.id]
	} else {
		log.Println("Ad", id, "not found")
	}
	return nil;
}
*/
func (r *Targets)LoadAll(db *sql.DB) {
	log.Print("Load targets.")
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM target WHERE !(status&128) LIMIT 3000000", TARGET_SQL_FIELDS))
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	r.list=make(TargetsList)
	for rows.Next() {
		var t=new(Target)
		t.setFromRow(rows)
		//log.Println(ad, len(r.ads))
		r.list[t.id] = t;
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(r.list), "targets");
}
