local wezterm = require "wezterm"
local config = {}

config.color_scheme_dirs = {
  os.getenv("HOME") .. "/.config/wezterm/themes/bearded-theme",
}

config.color_scheme = "Bearded Theme Monokai Metallian"

return config
