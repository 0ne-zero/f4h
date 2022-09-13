package controller_helper

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/database/model_function"
	viewmodel "github.com/0ne-zero/f4h/public_struct/view_model"
	general_func "github.com/0ne-zero/f4h/utilities/functions/general"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetTitle(t string) string {
	return constansts.AppName + fmt.Sprintf(" | %s", t)
}

func IsXMRWalletAddrValid(addr *string) bool {
	return false
}
func ManageWalletValidation(wallet_addr string) error {
	if wallet_addr == "" {
		return fmt.Errorf("Fill the wallet address field")
	}
	if !IsXMRWalletAddrValid(&wallet_addr) {
		return fmt.Errorf("The wallet address is invalid")
	}
	return nil
}
func ManageAddressValidation(addr *viewmodel.UserPanel_Profile_ManageAddress) error {
	if addr.Name == "" {
		return fmt.Errorf("Fill the name field")
	}
	if addr.Country == "" {
		return fmt.Errorf("Fill the countary field")
	}
	if addr.Province == "" {
		return fmt.Errorf("Fill the province field")
	}
	if addr.City == "" {
		return fmt.Errorf("Fill the city field")
	}
	if addr.Street == "" {
		return fmt.Errorf("Fill the street field")
	}
	if addr.BuildingNumber == "" {
		return fmt.Errorf("Fill the building number field")
	}
	if addr.PostalCode == "" {
		return fmt.Errorf("Fill the postal code field")
	}
	if addr.Description == "" {
		return fmt.Errorf("Fill the description field")
	}

	return nil
}
func EditAccountValidation(email, new_pass, new_pass_confirm, cur_pass, original_pass_hash *string) error {
	cur_pass_hash, err := general_func.Hashing(*cur_pass)
	if err != nil {
		return err
	}
	if cur_pass_hash != *original_pass_hash {
		return fmt.Errorf("Current password is incorrect")
	}
	if *email == "" {
		return fmt.Errorf("Fill the email field")
	}
	if *new_pass_confirm != "" && *new_pass == "" {
		return fmt.Errorf("Fill the new password field")
	}

	if *new_pass != "" {
		if *new_pass_confirm == "" {
			return fmt.Errorf("Fill the new password confirm field")
		}
		if len(*new_pass) != len(*new_pass_confirm) || *new_pass != *new_pass_confirm {
			return fmt.Errorf("The passwords isn't same")
		}
	}

	if *cur_pass == "" {
		return fmt.Errorf("Fill the current password field")
	}
	return nil
}

// Returns param value as string
// Returns "" if there is no such value with that key
func GetParamFromReferer(ref, key string) string {
	// Without parameter
	if !strings.Contains(ref, "?") {
		return ""
	}
	query := strings.Split(ref, "?")[1]
	if strings.Contains(query, "?") {
		query = strings.Replace(query, "?", "", 1)
	}
	raw_parameters := strings.Split(query, "&")
	var parameters map[string]string
	for i := range raw_parameters {
		param_parts := strings.SplitN(raw_parameters[i], "=", 1)
		param_name := param_parts[0]
		param_value := param_parts[1]
		parameters[param_name] = param_value
	}

	if value, ok := parameters[key]; ok {
		return value
	} else {
		return ""
	}

}
func MakeProfileViewData(c *gin.Context, use_referer bool) (gin.H, error) {
	// Available tabs and their modes
	tabs_modes := map[string][]string{
		"overview": {"front_page", "logins", "orders"},
		"profile":  {"edit_account", "edit_signature", "edit_avatar", "manage_login", "manage_address", "manage_wallets"},
		"products": {"front_page"},
		"payments": {"front_page"},
		"topics":   {"front_page"},
		"polls":    {"front_page"},
	}
	// Get user id
	user_id := sessions.Default(c).Get("UserID").(int)
	// Sort tabs and modes
	sorted_tabs := general_func.GetMapKeys(tabs_modes)
	sort.Strings(sorted_tabs)

	// Finally view data
	view_data := gin.H{
		"Title": fmt.Sprintf("%s Profile | %s", sessions.Default(c).Get("Username"), constansts.AppName),
		"Tabs":  sorted_tabs,
	}
	var tab string
	var mode string
	if use_referer {
		tab = GetParamFromReferer(c.Request.Referer(), "tab")
		mode = GetParamFromReferer(c.Request.Referer(), "mode")
	} else {
		tab = c.Query("tab")
		mode = c.Query("mode")
	}
	// Check is user selected tab
	if tab != "" {
		// Check selected tab is available
		if !general_func.ExistsStringInStringSlice(tab, general_func.GetMapKeys(tabs_modes)) {
			return nil, fmt.Errorf("Entered unavailable tab")
		}
		// Insert tab to information that will send to template
		view_data["Tab"] = tab
		// Check user selected tab mode
		if mode != "" {
			// Check entered mode available
			if !general_func.ExistsStringInStringSlice(mode, tabs_modes[tab]) {
				return nil, fmt.Errorf("Entered unavailable mode")
			}
			// Set tab modes
			view_data["TabModes"] = tabs_modes[tab]
			// Insert mode to information that will send to template
			view_data["Mode"] = mode
			// Get panel data
			switch tab {
			case "overview":
				switch mode {
				case "front_page":
					panel_data, err := model_function.GetUserDataForUserPanel_Overview_FrontPage(user_id)
					if err != nil {
						return nil, err
					}
					view_data["PanelData"] = panel_data
				case "logins":
					panel_data, err := model_function.GetUserDataForUserPanel_Overview_Logins(user_id)
					if err != nil {
						return nil, err
					}
					view_data["PanelData"] = panel_data
				case "orders":
					panel_data, err := model_function.GetUserDataForUserPanel_Overview_Orders(user_id)
					if err != nil {
						return nil, err
					}
					view_data["PanelData"] = panel_data
				}
			case "profile":
				switch mode {
				case "edit_account":
					panel_data, err := model_function.GetUserDataForUserPanel_Profile_EditAccount(user_id)
					if err != nil {
						return nil, err
					}
					view_data["PanelData"] = panel_data
				case "edit_signature":
					panel_data, err := model_function.GetUserDataForUserPanel_Profile_EditSignature(user_id)
					if err != nil {
						return nil, err
					}
					view_data["PanelData"] = panel_data
				case "edit_avatar":
					panel_data, err := model_function.GetUserDataForUserPanel_Profile_EditAvatar(user_id)
					if err != nil {
						return nil, err
					}
					view_data["PanelData"] = panel_data
				case "manage_address":
					panel_data, err := model_function.GetUserDataForUserPanel_Profile_ManageAddress(user_id)
					if err != nil {
						return nil, err
					}
					view_data["PanelData"] = panel_data
				case "manage_wallets":
					panel_data, err := model_function.GetUserDataForUserPanel_Profile_ManageWallet(user_id)
					if err != nil {
						return nil, err
					}
					view_data["PanelData"] = panel_data
				}
			case "products":
				switch mode {
				case "front_page":
				}
			case "payments":
				switch mode {
				case "front_page":
				}
			case "topics":
				switch mode {
				case "front_page":
				}
			case "polls":
				switch mode {
				case "front_page":
				}
			}

			return view_data, nil
		}
		// Set tab modes
		view_data["TabModes"] = tabs_modes[tab]
		// Tab's mode isn't selected so select default mode (first element)
		view_data["Mode"] = tabs_modes[tab][0]
		// Get panel data
		switch tab {
		case "overview":
			data, err := model_function.GetUserDataForUserPanel_Overview_FrontPage(user_id)
			if err != nil {
				return nil, err
			}
			view_data["PanelData"] = data
		case "profile":
			data, err := model_function.GetUserDataForUserPanel_Profile_EditAccount(user_id)
			if err != nil {
				return nil, err
			}
			view_data["PanelData"] = data
		case "products":
		case "payments":
		case "topics":
		case "polls":
		default:
		}
		return view_data, nil
	}

	// Neither tab nor tab mode is selected, so select default tab and default tab mode, which that means overview tab and its first mode
	view_data["Tab"] = "overview"
	view_data["Mode"] = "front_page"
	view_data["TabModes"] = tabs_modes["overview"]

	// Get panel data
	panel_data, err := model_function.GetUserDataForUserPanel_Overview_FrontPage(user_id)
	if err != nil {
		return nil, err
	}
	view_data["PanelData"] = panel_data
	return view_data, nil
}

// If returned value be nil means there is no error(data are validated)
func AddProductValidation(name, price, inventory, description, tags *string) gin.H {
	// Emptiness check
	if *name == "" {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Price":       price,
				"Inventory":   inventory,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Please fill the title field"},
		}
		return view_data
	}
	if *price == "" {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Inventory":   inventory,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Please fill the price field"},
		}
		return view_data
	}
	if *inventory == "" {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Price":       price,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Please fill the inventory field"},
		}
		return view_data
	}
	if *description == "" {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Price":     price,
				"Inventory": inventory,
				"Name":      name,
				"Tags":      tags,
			},
			"NotesData": gin.H{"Error": "Please fill the description"},
		}
		return view_data
	}
	if *tags == "" {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Price":       price,
				"Inventory":   inventory,
				"Description": description,
			},
			"NotesData": gin.H{"Error": "Please fill the tags field"},
		}
		return view_data
	}
	// Convert to actual type
	_, err := strconv.ParseFloat(*price, 64)
	if err != nil {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Inventory":   inventory,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Please fill the price field as number (in Dollor)"},
		}
		return view_data
	}
	_, err = strconv.Atoi(*inventory)
	if err != nil {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Price":       price,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Please fill the inventory field as number"},
		}
		return view_data
	}
	// Check length
	if len(*name) < 5 {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Price":       price,
				"Inventory":   inventory,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Name field is too short, It should be minimum 5 character"},
		}
		return view_data
	}
	if len(*description) < 50 {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Price":       price,
				"Inventory":   inventory,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Description field is too short, It should be minimum 50 character"},
		}
		return view_data
	}

	if strings.Contains(*tags, "|") {
		splitted_tags := strings.Split(*tags, "|")
		splitted_tags_len := len(splitted_tags)
		if splitted_tags_len > 5 {
			view_data := gin.H{
				"ProductInfo": gin.H{
					"Name":        name,
					"Price":       price,
					"Inventory":   inventory,
					"Description": description,
					"Tags":        tags,
				},
				"NotesData": gin.H{"Error": fmt.Sprintf("Maximum number of tag is 5, But you inserted %d tags", splitted_tags_len)},
			}
			return view_data
		}
	}
	return nil
}
func EditProductValidation(id, name, price, inventory, description, tags *string) gin.H {
	return AddProductValidation(name, price, inventory, description, tags)
}
func ErrorPage(c *gin.Context, user_msg string) {
	// Add bad request to database
	AddBadRequest(c)

	// Return response
	var view_data = gin.H{}
	view_data["Title"] = "Error"
	view_data["Error"] = user_msg
	c.HTML(http.StatusInternalServerError, "error.html", view_data)
}
func GetIDFromURLParameters(c *gin.Context) (int, error) {
	id_str := c.Param("id")
	if id_str == "" {
		return 0, fmt.Errorf("there isn't id in parameters")
	}
	id_int, err := strconv.Atoi(id_str)
	return id_int, err
}
func AddBadRequest(c *gin.Context) error {
	var bad_request = model.BadRequest{
		IP:     c.ClientIP(),
		Url:    c.Request.URL.Path,
		Method: c.Request.Method,
		Time:   time.Now().UTC(),
	}
	return model_function.Add(&bad_request)
}
func ClientInfoInMap(c *gin.Context) map[string]string {
	return map[string]string{
		"IP":     c.ClientIP(),
		"URL":    c.Request.URL.Path,
		"METHOD": c.Request.Method,
	}
}
func DeleteUserSession(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
}
func DeleteUserTopicDraftFromSession(c *gin.Context) {
	s := sessions.Default(c)
	s.Delete("TopicSubject")
	s.Delete("TopicMarkdown")
	s.Delete("TopicTags")
	s.Save()
}
