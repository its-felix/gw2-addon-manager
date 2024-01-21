local base = require("github.com/its-felix/gw2-addon-manager/base/standalone.lua")
-- local base = require("github.com/its-felix/gw2-addon-manager/base/addonloader.lua")
-- local base = require("github.com/its-felix/gw2-addon-manager/base/arcdps.lua")

return {
    name = "DRF",
    find_installed = function(api)
        local f = api.file("addons/d3d11.dll")
        if f ~= nil and api.sha256(f) == "abcdef" {
            return "v1.0.0"
        }

        return nil
    end,
    install = function(api)

    end,
    uninstall = function(api)

    end,
    arguments = {

    },
}