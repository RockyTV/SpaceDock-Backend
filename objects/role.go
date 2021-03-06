/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
    "encoding/json"
)

type Role struct {
    Model

    Name      string `gorm:"size:128;unique_index;not null" spacedock:"lock"`
    Params    string `gorm:"size:4096" spacedock:"json" spacedock:"lock"`
    Abilities []Ability `gorm:"many2many:role_abilities" json"-" spacedock:"lock"`
    Users     []User `gorm:"many2many:role_users" json"-" spacedock:"lock"`
}

func (s *Role) AfterFind() {
    app.DBRecursionLock.Lock()
    if _, ok := app.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        app.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if app.DBRecursion[utils.CurrentGoroutineID()] >= app.DBRecursionMax {
        app.DBRecursionLock.Unlock()
        return
    }
    isRoot := app.DBRecursion[utils.CurrentGoroutineID()] == 0
    app.DBRecursion[utils.CurrentGoroutineID()] += 1
    app.DBRecursionLock.Unlock()

    app.Database.Model(s).Related(&(s.Abilities), "Abilities")
    app.Database.Model(s).Related(&(s.Users), "Users")

    app.DBRecursionLock.Lock()
    app.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(app.DBRecursion, utils.CurrentGoroutineID())
    }
    app.DBRecursionLock.Unlock()
}

func (role *Role) AddAbility(name string) *Ability {
    ability := &Ability {}
    app.Database.Where("name = ?", name).First(ability)
    if ability.Name != name {
        ability.Name = name
        ability.Meta = "{}"
        app.Database.Save(ability)
    }
    role.Abilities = append(role.Abilities, *ability)
    app.Database.Save(role).Save(ability)
    return &role.Abilities[len(role.Abilities) - 1]
}

func (role *Role) RemoveAbility(name string) {
    ability := &Ability {}
    app.Database.Where("name = ?", name).First(ability)
    if ability.Name == "" {
        return
    }
    if e,i := utils.ArrayContains(ability, role.Abilities); e {
        role.Abilities = append(role.Abilities[:i], role.Abilities[i + 1:]...)
        app.Database.Save(role)
    }
}

func (role *Role) HasAbility(name string) bool {
    ability := &Ability {}
    app.Database.Where("name = ?", name).First(ability)
    if ability.Name == "" {
        return false
    }
    e,_ := utils.ArrayContains(ability, &(role.Abilities))
    return e
}

func (role *Role) GetParams(ability string, param string) []string {
    var temp map[string]map[string][]string
    err := json.Unmarshal([]byte(role.Params), &temp)
    if err != nil {
        return nil
    }
    if _,ok := temp[ability]; ok {
        if _,ok := temp[ability][param]; ok {
            return temp[ability][param]
        }
    }
    return nil
}

func (role *Role) AddParam(ability string, param string, value string) error  {
    var temp map[string]map[string][]string
    err := json.Unmarshal([]byte(role.Params), &temp)
    if err != nil {
        return err
    }
    if _,ok := temp[ability]; !ok {
        temp[ability] = map[string][]string{}
    }
    if _,ok := temp[ability][param]; !ok {
        temp[ability][param] = []string{}
    }
    if ok,_ := utils.ArrayContains(value, temp[ability][param]); !ok {
        temp[ability][param] = append(temp[ability][param], value)
    }
    val,err := json.Marshal(temp)
    if err != nil {
        return err
    }
    role.Params = string(val)
    return nil
}

func (role *Role) RemoveParam(ability string, param string, value string) error {
    var temp map[string]map[string][]string
    err := json.Unmarshal([]byte(role.Params), &temp)
    if err != nil {
        return err
    }
    if _,ok := temp[ability]; ok {
        if _,ok := temp[ability][param]; ok {
            if ok,index := utils.ArrayContains(value, temp[ability][param]); ok {
                temp[ability][param] = append(temp[ability][param][:index], temp[ability][param][index+1:]...)
            }
        }
    }
    val,err := json.Marshal(temp)
    if err != nil {
        return err
    }
    role.Params = string(val)
    return nil
}