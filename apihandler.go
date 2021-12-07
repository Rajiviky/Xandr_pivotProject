package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"pivot/src/config"

	"pivot/pkg/github.com/gorilla/mux"
)

type EmpSkill struct {
	Skill_nm    string
	Proficiency string
	Version     string
	IsPrimary   string
	Last_used   string
	Total_exp   int
	Attuid      string
	First_nm    string
	Last_nm     string
	Skill_id    int
}

type Emp struct {
	Attuid   string
	First_nm string
	Last_nm  string
	Email    string
	Mgr_id   string
	Status   string
}

type SkillNew struct {
	Skill_nm string
}

var Eskill []EmpSkill

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello World from apiHandler!</h1>")
}

// Get the employee skill details
func GetEmployeeskill(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	fmt.Println("Getting employee detail")
	//User api handling
	m := map[string]string{}
	var apiRes Response
	res, err := os.Open("./static/Employee.json")
	fmt.Println("API call comleted")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//defer res.Body.Close()
	defer res.Close()

	data, _ := ioutil.ReadAll(res)
	//fmt.Println("data=", string(data))

	json.Unmarshal(data, &apiRes)
	//fmt.Println("Result=", apiRes)
	var person_id int
	for _, i := range apiRes.Results {
		if i.Login == params["attuid"] {
			fname := i.First_name
			lname := i.Last_name
			person_id = i.Person_id

			m[i.Login] = strings.Join([]string{fname, lname}, "-")

			break
		}
	}
	// check if the user is a manager/supervisor

	for _, i := range apiRes.Results {
		if i.Manager_person_id == person_id {
			fname := i.First_name
			lname := i.Last_name

			m[i.Login] = strings.Join([]string{fname, lname}, "-")

		}

	}
	//params["attuid"]
	var users = "'" + params["attuid"] + "'"
	for k := range m {
		//users = "'" + users + "'" + "," + "'" + k + "'"
		users = users + "," + "'" + k + "'"

	}
	//fmt.Println("users=", users)
	// done
	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Failed to connect to database")
	}
	var Skill EmpSkill
	Eskill = nil

	stmt := `SELECT b.skill_nm,
	CASE WHEN a.proficiency = 1 THEN "Beginner"
	WHEN a.proficiency = 2 THEN "Intermediate"
	WHEN a.proficiency = 3 THEN "Advance"
	WHEN a.proficiency = 4 THEN "Expert"
	END as proficiency
	,a.version,a.is_primary,a.total_exp_in_month ,strftime('%Y-%m-%d',a.last_used),a.attuid ,b.skill_id
	FROM employeeskill a INNER JOIN skill b ON a.skill_id=b.skill_id where a.attuid in (` + users + ") order by case when a.attuid='" + params["attuid"] + "' then 1 else 2 end,a.attuid;"
	rows, err1 := db.Query(stmt)
	//fmt.Println("stmt=", stmt)
	if err1 != nil {
		fmt.Println("Failed to fetch the employee detail of user.")
		fmt.Println(err1)
	} else {
		for rows.Next() {
			err2 := rows.Scan(&Skill.Skill_nm, &Skill.Proficiency, &Skill.Version, &Skill.IsPrimary, &Skill.Total_exp, &Skill.Last_used, &Skill.Attuid, &Skill.Skill_id)
			if err2 != nil {
				fmt.Println("Error in fetching the skill details")
			}

			Skill.First_nm = strings.Split(m[Skill.Attuid], "-")[0]
			Skill.Last_nm = strings.Split(m[Skill.Attuid], "-")[1]

			Eskill = append(Eskill, Skill)

		}
	}
	fmt.Println(Eskill)
	json.NewEncoder(w).Encode(Eskill)
	defer rows.Close()

}

// GET all employees
func GetEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Emp
	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Failed to connect to database")
	}
	//var Skill EmpSkill

	stmt := "SELECT attuid,first_nm,last_nm,email,mgr_id,status FROM employee;"
	rows, err1 := db.Query(stmt)
	if err1 != nil {
		fmt.Println("Failed to fetch the employee detail of user.")
		fmt.Println(err1)
	} else {
		for rows.Next() {
			err2 := rows.Scan(&emp.Attuid, &emp.First_nm, &emp.Last_nm, &emp.Email, &emp.Mgr_id, &emp.Status)
			if err2 != nil {
				fmt.Println("Error in fetching the skill details")
			}
			json.NewEncoder(w).Encode(emp)

		}
	}
	defer rows.Close()

}

//Add new employee
func AddEmployee(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	//param := mux.Vars(r)
	//fmt.Fprint(w, decoder)
	var e Emp
	err := decoder.Decode(&e)
	if err != nil {
		panic(err)
	}

	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Failed to connect to database")
	}

	stmt, _ := db.Prepare("insert into employee('attuid','first_nm','last_nm','email','mgr_id','status') values(?,?,?,?,?,?);")
	_, err1 := stmt.Exec(e.Attuid, e.First_nm, e.Last_nm, e.Email, e.Mgr_id, e.Status)
	if err1 != nil {
		fmt.Println("Failed to insert the new user.")
		fmt.Println(err1)
	} else {
		fmt.Println("Record added successfully")
	}

}

//Delete Employee Skill
func DeleteEmployeeSkill(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	attuid := strings.Split(param["attuid-skillid"], "-")[0]
	skillid := strings.Split(param["attuid-skillid"], "-")[1]

	fmt.Println(skillid, attuid)
	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Failed to connect to database")
	}

	stmt, _ := db.Prepare("DELETE FROM employeeskill where attuid=? and skill_id=?")
	_, err1 := stmt.Exec(attuid, skillid)
	if err1 != nil {
		fmt.Println("Failed to delete the record.")
		fmt.Println(err1)
	} else {
		fmt.Println("Record deleted successfully")
	}

}

//DELETE Employee
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Failed to connect to database")
	}

	stmt, _ := db.Prepare("DELETE FROM employee where attuid=?")
	_, err1 := stmt.Exec(param["attuid"])
	if err1 != nil {
		fmt.Println("Failed to delete the record.")
		fmt.Println(err1)
	} else {
		fmt.Println("Record deleted successfully")
	}

}

//Add new skill
func AddSkill(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var sk SkillNew
	err := decoder.Decode(&sk)
	if err != nil {
		panic(err)
	}

	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Failed to connect to database")
	}

	stmt, _ := db.Prepare("insert into skill('skill_nm') values(?);")
	_, err1 := stmt.Exec(sk.Skill_nm)
	if err1 != nil {
		fmt.Println("Failed to insert new skill.")
		fmt.Println(err1)
	} else {
		fmt.Println("Skill added successfully")
	}

}
