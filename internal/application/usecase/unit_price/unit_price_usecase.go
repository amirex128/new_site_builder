package unitpriceusecase

import (
	"fmt"
	unit_price2 "github.com/amirex128/new_site_builder/internal/application/dto/unit_price"
	"github.com/amirex128/new_site_builder/internal/application/usecase"
	"github.com/amirex128/new_site_builder/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/internal/contract"
	"github.com/amirex128/new_site_builder/internal/contract/common"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	"github.com/amirex128/new_site_builder/internal/domain/enums"
)

type UnitPriceUsecase struct {
	*usecase.BaseUsecase
	unitPriceRepo repository2.IUnitPriceRepository
	userRepo      repository2.IUserRepository
}

func NewUnitPriceUsecase(c contract.IContainer) *UnitPriceUsecase {
	return &UnitPriceUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		unitPriceRepo: c.GetUnitPriceRepo(),
		userRepo:      c.GetUserRepo(),
	}
}

func (u *UnitPriceUsecase) UpdateUnitPriceCommand(params *unit_price2.UpdateUnitPriceCommand) (*resp.Response, error) {
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "only admins can update unit prices")
	}
	existingUnitPrice, err := u.unitPriceRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	if params.Name != nil {
		existingUnitPrice.Name = *params.Name
	}
	if params.Price != nil {
		existingUnitPrice.Price = int64(*params.Price)
	}
	if params.HasDay != nil {
		existingUnitPrice.HasDay = *params.HasDay
	}
	if params.DiscountType != nil {
		existingUnitPrice.DiscountType = *params.DiscountType
	}
	if params.Discount != nil {
		discount := int64(*params.Discount)
		existingUnitPrice.Discount = &discount
	}
	err = u.unitPriceRepo.Update(existingUnitPrice)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Updated, map[string]interface{}{"unitPrice": existingUnitPrice}, "Unit price updated successfully"), nil
}

func (u *UnitPriceUsecase) CalculateUnitPriceQuery(params *unit_price2.CalculateUnitPriceQuery) (*resp.Response, error) {
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "user not authenticated")
	}
	var unitPriceNames []string
	for _, up := range params.UnitPrices {
		unitPriceNames = append(unitPriceNames, string(*up.UnitPriceName))
	}
	allUnitPricesResult, err := u.unitPriceRepo.GetAll(common.PaginationRequestDto{Page: 1, PageSize: 100})
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	currentUser, err := u.userRepo.GetByID(*userID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	_ = currentUser
	var data []map[string]interface{}
	for _, unitPriceParam := range params.UnitPrices {
		var matchingUnitPrice *enums.UnitPriceNameEnum
		for _, up := range allUnitPricesResult.Items {
			if up.Name == string(*unitPriceParam.UnitPriceName) {
				matchingUnitPrice = unitPriceParam.UnitPriceName
				if *unitPriceParam.UnitPriceName == "storage_mb_credits" && unitPriceParam.UnitPriceDay != nil {
					finalPrice := up.Price * int64(*unitPriceParam.UnitPriceCount) * int64(*unitPriceParam.UnitPriceDay)
					discountAmount := int64(0)
					if up.Discount != nil && *up.Discount > 0 {
						if up.DiscountType == "fixed" {
							discountAmount = *up.Discount
							finalPrice = finalPrice - discountAmount
						} else if up.DiscountType == "percentage" {
							discountAmount = (finalPrice * (*up.Discount)) / 100
							finalPrice = finalPrice - discountAmount
						}
					}
					if finalPrice < 0 {
						finalPrice = 0
					}
					data = append(data, map[string]interface{}{
						"UnitPriceName":  *unitPriceParam.UnitPriceName,
						"UnitPriceCount": *unitPriceParam.UnitPriceCount,
						"TotalPrice":     finalPrice + discountAmount,
						"FinalPrice":     finalPrice,
						"DiscountAmount": discountAmount,
					})
				} else {
					finalPrice := up.Price * int64(*unitPriceParam.UnitPriceCount)
					discountAmount := int64(0)
					if up.Discount != nil && *up.Discount > 0 {
						if up.DiscountType == "fixed" {
							discountAmount = *up.Discount
							finalPrice = finalPrice - discountAmount
						} else if up.DiscountType == "percentage" {
							discountAmount = (finalPrice * (*up.Discount)) / 100
							finalPrice = finalPrice - discountAmount
						}
					}
					if finalPrice < 0 {
						finalPrice = 0
					}
					data = append(data, map[string]interface{}{
						"UnitPriceName":  *unitPriceParam.UnitPriceName,
						"UnitPriceCount": *unitPriceParam.UnitPriceCount,
						"TotalPrice":     finalPrice + discountAmount,
						"FinalPrice":     finalPrice,
						"DiscountAmount": discountAmount,
					})
				}
				break
			}
		}
		if matchingUnitPrice == nil {
			return nil, resp.NewError(resp.NotFound, fmt.Sprintf("unit price not found for name: %s", string(*unitPriceParam.UnitPriceName)))
		}
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{"calculatedPrices": data}, "Unit prices calculated successfully"), nil
}

func (u *UnitPriceUsecase) GetAllUnitPriceQuery(params *unit_price2.GetAllUnitPriceQuery) (*resp.Response, error) {
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "only admins can get all unit prices")
	}
	unitPricesResult, err := u.unitPriceRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     unitPricesResult.Items,
		"total":     unitPricesResult.TotalCount,
		"page":      unitPricesResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": unitPricesResult.TotalPages,
	}, "Unit prices retrieved successfully"), nil
}
