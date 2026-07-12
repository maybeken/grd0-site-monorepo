package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"

	"grd0.net/api/schema"
)

func (h *Handler) GetCoffeeBeans(c echo.Context) error {
	db := h.DB

	var beans []schema.CoffeeBean
	db.Order("created_at DESC").Find(&beans)

	return c.JSON(http.StatusOK, beans)
}

func (h *Handler) UpsertCoffeeBean(c echo.Context) error {
	db := h.DB

	var bean schema.CoffeeBean
	if err := c.Bind(&bean); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	if bean.Name == "" {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Bean name is required.")
	}

	if bean.ID != uuid.Nil {
		res := db.Model(&schema.CoffeeBean{}).Where("id = ?", bean.ID).Select("*").Updates(&bean)
		if res.Error != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, res.Error.Error())
		}
		if res.RowsAffected == 0 {
			return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
		}
	} else {
		if err := db.Create(&bean).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, bean)
}

func (h *Handler) DeleteCoffeeBean(c echo.Context) error {
	id := c.Param("id")

	db := h.DB

	res := db.Where("id = ?", id).Delete(&schema.CoffeeBean{})
	if res.Error != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, res.Error.Error())
	} else if res.RowsAffected <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusAccepted, nil)
}

func (h *Handler) GetCoffeeEquipment(c echo.Context) error {
	db := h.DB

	var equipment []schema.BrewEquipment
	db.Order("created_at DESC").Find(&equipment)

	return c.JSON(http.StatusOK, equipment)
}

func (h *Handler) UpsertCoffeeEquipment(c echo.Context) error {
	db := h.DB

	var equipment schema.BrewEquipment
	if err := c.Bind(&equipment); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	if equipment.Name == "" {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Equipment name is required.")
	}

	if equipment.ID != uuid.Nil {
		res := db.Model(&schema.BrewEquipment{}).Where("id = ?", equipment.ID).Select("*").Updates(&equipment)
		if res.Error != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, res.Error.Error())
		}
		if res.RowsAffected == 0 {
			return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
		}
	} else {
		if err := db.Create(&equipment).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, equipment)
}

func (h *Handler) DeleteCoffeeEquipment(c echo.Context) error {
	id := c.Param("id")

	db := h.DB

	res := db.Where("id = ?", id).Delete(&schema.BrewEquipment{})
	if res.Error != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, res.Error.Error())
	} else if res.RowsAffected <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusAccepted, nil)
}

func (h *Handler) GetCoffeeTastings(c echo.Context) error {
	db := h.DB

	var tastings []schema.TastingNote
	db.Preload("Bean").Preload("Equipment").Preload("Grinder").Order("pinned DESC").Order("tasted_at DESC").Find(&tastings)

	return c.JSON(http.StatusOK, tastings)
}

func (h *Handler) UpsertCoffeeTasting(c echo.Context) error {
	db := h.DB

	var tasting schema.TastingNote
	if err := c.Bind(&tasting); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	tasting.BeanID = tasting.Bean.ID
	tasting.EquipmentID = tasting.Equipment.ID
	tasting.GrinderID = tasting.Grinder.ID

	if tasting.TastedAt.IsZero() {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Tasted date is required.")
	}

	if tasting.GrinderID != uuid.Nil {
		var equipment schema.BrewEquipment
		db.First(&equipment, "id = ?", tasting.GrinderID)

		if equipment.Type == nil || *equipment.Type != "Grinder" {
			return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Grinder ID provided is not a grinder.")
		}
	}

	if tasting.ID != uuid.Nil {
		res := db.Omit(clause.Associations).Model(&schema.TastingNote{}).Where("id = ?", tasting.ID).Select("*").Updates(&tasting)
		if res.Error != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, res.Error.Error())
		}
		if res.RowsAffected == 0 {
			return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
		}
	} else {
		if err := db.Omit(clause.Associations).Create(&tasting).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	}

	db.Preload("Bean").Preload("Equipment").Preload("Grinder").First(&tasting, "id = ?", tasting.ID)

	return c.JSON(http.StatusOK, tasting)
}

func (h *Handler) DeleteCoffeeTasting(c echo.Context) error {
	id := c.Param("id")

	db := h.DB

	res := db.Where("id = ?", id).Delete(&schema.TastingNote{})
	if res.Error != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, res.Error.Error())
	} else if res.RowsAffected <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusAccepted, nil)
}
