/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package routes

import (
    "gopkg.in/kataras/iris.v6"
    "SpaceDock/objects"
    "SpaceDock"
    "SpaceDock/utils"
    "github.com/spf13/cast"
    "SpaceDock/middleware"
)

/*
 Registers the routes for the featured section
 */
func FeaturedRegister() {
    Register(GET, "/api/mods/featured", list_featured)
    Register(GET, "/api/mods/featured/:gameshort", list_featured_game)
    Register(POST, "/api/mods/featured/:gameshort",
        middleware.NeedsPermission("mods-feature", true, "gameshort"),
        add_featured,
    )
    Register(DELETE, "/api/mods/featured/:gameshort",
        middleware.NeedsPermission("mods-feature", true, "gameshort"),
        remove_featured,
    )
}

/*
 Path: /api/mods/featured
 Method: GET
 Description: Returns a list of featured mods.
 */
func list_featured(ctx *iris.Context) {
    var featured []objects.Featured
    SpaceDock.Database.Find(&featured)
    output := make([]map[string]interface{}, len(featured))
    for i,element := range featured {
        output[i] = utils.ToMap(element)
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(output), "data": output})
}

/*
 Path: /api/mods/featured/:gameshort
 Method: GET
 Description: Returns a list of featured mods for a specific game.
 */
func list_featured_game(ctx *iris.Context) {
    gameshort := ctx.GetString("gameshort")

    // Check if the game exists
    game := &objects.Game{}
    SpaceDock.Database.Where("short = ?", gameshort).First(game)
    if game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }

    var featured []objects.Featured
    SpaceDock.Database.Find(&featured)
    output := []map[string]interface{}{}
    for _,element := range featured {
        if element.Mod.Game.Short == gameshort {
            output = append(output, utils.ToMap(element))
        }
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(output), "data": output})
}

/*
 Path: /api/mods/featured/:gameshort
 Method: POST
 Description: Features a mod for this game. Required fields: modid
 Abilities: mods-feature
 */
func add_featured(ctx *iris.Context) {
    // Get the mod
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(utils.GetJSON(ctx, "modid"))

    // Get the mod
    mod := &objects.Mod{}
    SpaceDock.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid.").Code(2130))
        return
    } else if mod.Game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    } else if !mod.Published {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The mod must be published first.").Code(3022))
        return
    }

    // Check if the mod is already featured
    feature := &objects.Featured{}
    SpaceDock.Database.Where("mod_id = ?", modid).First(feature)
    if feature.ModID == modid {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The mod is already featured").Code(3015))
        return
    }

    // Everything is fine, lets feature the mod
    feature = objects.NewFeatured(*mod)
    SpaceDock.Database.Save(feature)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(feature)})
}

/*
 Path: /api/mods/featured/:gameshort
 Method: DELETE
 Description: Unfeatures a mod for this game. Required fields: modid
 Abilities: mods-feature
 */
func remove_featured(ctx *iris.Context) {
    // Get the mod
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(utils.GetJSON(ctx, "modid"))

    // Get the mod
    mod := &objects.Mod{}
    SpaceDock.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid.").Code(2130))
        return
    } else if mod.Game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    } else if !mod.Published {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The mod must be published first.").Code(3022))
        return
    }

    // Check if the mod is already featured
    feature := &objects.Featured{}
    SpaceDock.Database.Where("mod_id = ?", modid).First(feature)
    if feature.ModID != modid {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("This mod isn't featured.").Code(3015))
        return
    }

    // Everything is fine, lets remove the feature
    SpaceDock.Database.Delete(feature)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}