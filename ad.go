package main

import (
	"fmt"
	//"io/ioutil"
	"log"
	"database/sql"
	"sync"
)

//type extType

type Ad struct {
	id intc
	campaign intc
	target Target
	status uint16
	width uint16
	height uint16
	ext uint8
	text string
	text1 string
	url string
	price intc
	sync sync.Mutex
}

type AdsList map[intc]*Ad

type Ads struct {
	list AdsList
}

const AD_SQL = "SELECT ad.id, campaignID, targetID, ad.status, width, height, ext, text, text1, url, price, t.status AS tstatus, priority, aclHours, aclWeekDay, aclRegion, aclAgeSex, aclCats, retargetID, viewsUserDayLimit, periods, viewsDayLimit, viewsToday, viewsYesterday, views, clicksToday, clicksYesterday, clicks, spentLimit, spent, spentDayLimit, spentToday, maxPrice FROM ad JOIN target AS t ON t.id=ad.targetID WHERE !(ad.status&128) AND t.status&1"

func (ad *Ad)setFromRow(rows *sql.Rows) {
	target:=ad.target
	periods:=""
	err := rows.Scan(&ad.id, &ad.campaign, &target.id, &ad.status, &ad.width, &ad.height, &ad.ext, &ad.text, &ad.text1, &ad.url, &ad.price,
			&target.status, &target.priority, &target.aclHours, &target.aclWeekDay, &target.aclRegion, &target.aclAgeSex, &target.aclCats, &target.retargetID, &target.viewsUserDayLimit, &periods, &target.viewsDayLimit , &target.viewsToday, &target.viewsYesterday, &target.views, &target.clicksToday, &target.clicksYesterday, &target.clicks, &target.spentLimit, &target.spent, &target.spentDayLimit, &target.spentToday, &target.maxPrice)
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

func (a *Ads)Load(id intc, db *sql.DB) *Ad {
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

//func (ad *Ad)Update(id intc) bool {
//	ad = ads.Load(id)
//}

func (ad *Ad)getJson() string {
	return fmt.Sprintf("{id:%d, img: \"%s\", text:\"%s\", text1:\"%s\", price:\"%d грн\"}", ad.id, ad.Path(), ad.text, ad.text1, ad.price)
}

func (ad *Ad)Path() string {
	return fmt.Sprintf("http://i.mediatraffic.com.ua/%d/%d/%d/%d.%s", 100, (ad.id %100 / 10), (ad.id % 10), ad.id, "jpg");
}
func (r *Ads)LoadAll(db *sql.DB) {
	log.Print("Load ads.")
	rows, err := db.Query(fmt.Sprintf("%s LIMIT 300000", AD_SQL))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	//log.Print(rows.affectedRows());
	r.list=make(AdsList)

	for rows.Next() {
		var ad=new(Ad)
		ad.setFromRow(rows)
		//log.Println(ad, len(r.ads))
		r.list[ad.id] = ad;
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(r.list), "ads");
}



