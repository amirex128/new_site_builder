package config

import (
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerDeliveryServiceVendorsListFallbackMode string `env:"SERVER_DELIVERY_SERVICE_VENDORS_LIST_FALLBACK_MODE"`
	ServerLogPath                                string `env:"SERVER_LOG_PATH"`
	ServerPort                                   string `env:"SERVER_PORT"`
	ServerLogToFile                              string `env:"SERVER_LOG_TO_FILE"`
	ServerDevelopMode                            string `env:"SERVER_DEVELOP_MODE"`
	ServerSentryAllowed                          string `env:"SERVER_SENTRY_ALLOWED"`
	ServerSentryDsn                              string `env:"SERVER_SENTRY_DSN"`

	ServerReadDbUser     string `env:"SERVER_READ_DB_USER"`
	ServerReadDbPassword string `env:"SERVER_READ_DB_PASSWORD"`
	ServerReadDbHost     string `env:"SERVER_READ_DB_HOST"`
	ServerReadDbPort     string `env:"SERVER_READ_DB_PORT"`
	ServerReadDbDb       string `env:"SERVER_READ_DB_DB"`

	ServerWriteDbUser     string `env:"SERVER_WRITE_DB_USER"`
	ServerWriteDbPassword string `env:"SERVER_WRITE_DB_PASSWORD"`
	ServerWriteDbHost     string `env:"SERVER_WRITE_DB_HOST"`
	ServerWriteDbPort     string `env:"SERVER_WRITE_DB_PORT"`
	ServerWriteDbDb       string `env:"SERVER_WRITE_DB_DB"`

	ServerSearchReadDbUser     string `env:"SERVER_SEARCH_READ_DB_USER"`
	ServerSearchReadDbPassword string `env:"SERVER_SEARCH_READ_DB_PASSWORD"`
	ServerSearchReadDbHost     string `env:"SERVER_SEARCH_READ_DB_HOST"`
	ServerSearchReadDbPort     string `env:"SERVER_SEARCH_READ_DB_PORT"`
	ServerSearchReadDbDb       string `env:"SERVER_SEARCH_READ_DB_DB"`

	ServerSearchRedisHost     string `env:"SERVER_SEARCH_REDIS_HOST"`
	ServerSearchRedisPassword string `env:"SERVER_SEARCH_REDIS_PASSWORD"`
	ServerSearchRedisDb       string `env:"SERVER_SEARCH_REDIS_DB"`

	ServerNewhomeRedisHost     string `env:"SERVER_NEWHOME_REDIS_HOST"`
	ServerNewhomeRedisPassword string `env:"SERVER_NEWHOME_REDIS_PASSWORD"`
	ServerNewhomeRedisDb       string `env:"SERVER_NEWHOME_REDIS_DB"`

	ServerFoodproRedisHost     string `env:"SERVER_FOODPRO_REDIS_HOST"`
	ServerFoodproRedisPassword string `env:"SERVER_FOODPRO_REDIS_PASSWORD"`
	ServerFoodproRedisDb       string `env:"SERVER_FOODPRO_REDIS_DB"`

	ServerPartyRedisHost     string `env:"SERVER_PARTY_REDIS_HOST"`
	ServerPartyRedisPassword string `env:"SERVER_PARTY_REDIS_PASSWORD"`
	ServerPartyRedisDb       string `env:"SERVER_PARTY_REDIS_DB"`

	ServerOsrmRedisHost     string `env:"SERVER_OSRM_REDIS_HOST"`
	ServerOsrmRedisPassword string `env:"SERVER_OSRM_REDIS_PASSWORD"`
	ServerOsrmRedisDb       string `env:"SERVER_OSRM_REDIS_DB"`

	ServerStockRedisHost     string `env:"SERVER_STOCK_REDIS_HOST"`
	ServerStockRedisPassword string `env:"SERVER_STOCK_REDIS_PASSWORD"`
	ServerStockRedisDb       string `env:"SERVER_STOCK_REDIS_DB"`

	ServerStockReserveRedisHost     string `env:"SERVER_STOCK_RESERVE_REDIS_HOST"`
	ServerStockReserveRedisPassword string `env:"SERVER_STOCK_RESERVE_REDIS_PASSWORD"`
	ServerStockReserveRedisDb       string `env:"SERVER_STOCK_RESERVE_REDIS_DB"`

	ServerBaseElasticHost        string `env:"SERVER_BASE_ELASTIC_HOST"`
	ServerBaseElasticPort        string `env:"SERVER_BASE_ELASTIC_PORT"`
	ServerBaseElasticInsecureTls string `env:"SERVER_BASE_ELASTIC_INSECURE_TLS"`
	ServerBaseElasticUsername    string `env:"SERVER_BASE_ELASTIC_USERNAME"`
	ServerBaseElasticPassword    string `env:"SERVER_BASE_ELASTIC_PASSWORD"`

	ServerNewElasticHost        string `env:"SERVER_NEW_ELASTIC_HOST"`
	ServerNewElasticPort        string `env:"SERVER_NEW_ELASTIC_PORT"`
	ServerNewElasticInsecureTls string `env:"SERVER_NEW_ELASTIC_INSECURE_TLS"`
	ServerNewElasticUsername    string `env:"SERVER_NEW_ELASTIC_USERNAME"`
	ServerNewElasticPassword    string `env:"SERVER_NEW_ELASTIC_PASSWORD"`

	ServerHomepageDistanceSort       string `env:"SERVER_HOMEPAGE_DISTANCE_SORT"`
	ServerBulkTreadCount             string `env:"SERVER_BULK_TREAD_COUNT"`
	ServerVendorType                 string `env:"SERVER_VENDOR_TYPE"`
	ServerSliderType                 string `env:"SERVER_SLIDER_TYPE"`
	ServerSearchVendorIndex          string `env:"SERVER_SEARCH_VENDOR_INDEX"`
	ServerSearchDealVendorIndex      string `env:"SERVER_SEARCH_DEAL_VENDOR_INDEX"`
	ServerSearchVendorType           string `env:"SERVER_SEARCH_VENDOR_TYPE"`
	ServerSearchSliderIndex          string `env:"SERVER_SEARCH_SLIDER_INDEX"`
	ServerSearchSliderType           string `env:"SERVER_SEARCH_SLIDER_TYPE"`
	ServerSearchProductIndex         string `env:"SERVER_SEARCH_PRODUCT_INDEX"`
	ServerSearchProductNewIndex      string `env:"SERVER_SEARCH_PRODUCT_NEW_INDEX"`
	ServerSearchVendorNewIndex       string `env:"SERVER_SEARCH_VENDOR_NEW_INDEX"`
	ServerSearchSliderNewIndex       string `env:"SERVER_SEARCH_SLIDER_NEW_INDEX"`
	ServerSearchProductType          string `env:"SERVER_SEARCH_PRODUCT_TYPE"`
	ServerAuthIntrospectUrl          string `env:"SERVER_AUTH_INTROSPECT_URL"`
	ServerAuthSessionUrl             string `env:"SERVER_AUTH_SESSION_URL"`
	ServerAuthSessionDomain          string `env:"SERVER_AUTH_SESSION_DOMAIN"`
	ServerResourceId                 string `env:"SERVER_RESOURCE_ID"`
	ServerResourceSecret             string `env:"SERVER_RESOURCE_SECRET"`
	ServerRamadanSahar               string `env:"SERVER_RAMADAN_SAHAR"`
	ServerRamadanEftar               string `env:"SERVER_RAMADAN_EFTAR"`
	ServerRamadanRespect             string `env:"SERVER_RAMADAN_RESPECT"`
	ServerRamadanAvini               string `env:"SERVER_RAMADAN_AVINI"`
	ServerRamadanStartDate           string `env:"SERVER_RAMADAN_START_DATE"`
	ServerRamadanEndDate             string `env:"SERVER_RAMADAN_END_DATE"`
	ServerRamadanVendorImage         string `env:"SERVER_RAMADAN_VENDOR_IMAGE"`
	ServerRamadanProductImage        string `env:"SERVER_RAMADAN_PRODUCT_IMAGE"`
	ServerRamadanSabaVendorImage     string `env:"SERVER_RAMADAN_SABA_VENDOR_IMAGE"`
	ServerRamadanSabaCategoryImage   string `env:"SERVER_RAMADAN_SABA_CATEGORY_IMAGE"`
	ServerRamadanSabaSupertypeImage  string `env:"SERVER_RAMADAN_SABA_SUPERTYPE_IMAGE"`
	ServerRamadanSabaProductImage    string `env:"SERVER_RAMADAN_SABA_PRODUCT_IMAGE"`
	ServerApiDomain                  string `env:"SERVER_API_DOMAIN"`
	ServerDeeplinkDomain             string `env:"SERVER_DEEPLINK_DOMAIN"`
	ServerShopUrl                    string `env:"SERVER_SHOP_URL"`
	ServerSearchOrderIndex           string `env:"SERVER_SEARCH_ORDER_INDEX"`
	ServerSearchOrderType            string `env:"SERVER_SEARCH_ORDER_TYPE"`
	ServerScheme                     string `env:"SERVER_SCHEME"`
	ServerNewhomeSectionCount        string `env:"SERVER_NEWHOME_SECTION_COUNT"`
	ServerNewhomeDesktopSectionCount string `env:"SERVER_NEWHOME_DESKTOP_SECTION_COUNT"`
	ServerSpecialLogToFile           string `env:"SERVER_SPECIAL_LOG_TO_FILE"`
	ServerUserSpecificSort           string `env:"SERVER_USER_SPECIFIC_SORT"`
	ServerGroundDistSupport          string `env:"SERVER_GROUND_DIST_SUPPORT"`
	ServerReplaceCatastrophe         string `env:"SERVER_REPLACE_CATASTROPHE"`
	ServerHardCatastrophe            string `env:"SERVER_HARD_CATASTROPHE"`

	ServerHappyHourStartHour        string `env:"SERVER_HAPPY_HOUR_START_HOUR"`
	ServerHappyHourEndHour          string `env:"SERVER_HAPPY_HOUR_END_HOUR"`
	ServerHappyHourDuration         string `env:"SERVER_HAPPY_HOUR_DURATION"`
	ServerHappyHourActiveDuration   string `env:"SERVER_HAPPY_HOUR_ACTIVE_DURATION"`
	ServerHappyHourInactiveDuration string `env:"SERVER_HAPPY_HOUR_INACTIVE_DURATION"`
	ServerGroundDistDeliverySupport string `env:"SERVER_GROUND_DIST_DELIVERY_SUPPORT"`
	ServerSecretToken               string `env:"SERVER_SECRET_TOKEN"`
	ServerGeoPrecision              string `env:"SERVER_GEO_PRECISION"`

	ServerNewCuisine                              string `env:"SERVER_NEW_CUISINE"`
	ElasticApmServerUrl                           string `env:"ELASTIC_APM_SERVER_URL"`
	ElasticApmServiceName                         string `env:"ELASTIC_APM_SERVICE_NAME"`
	ElasticApmSecretToken                         string `env:"ELASTIC_APM_SECRET_TOKEN"`
	ElasticApmActive                              string `env:"ELASTIC_APM_ACTIVE"`
	ElasticApmServiceNodeName                     string `env:"ELASTIC_APM_SERVICE_NODE_NAME"`
	ServerGrpcUrl                                 string `env:"SERVER_GRPC_URL"`
	ServerGrpcActive                              string `env:"SERVER_GRPC_ACTIVE"`
	ServerGrpcFoodStoryActive                     string `env:"SERVER_GRPC_FOOD_STORY_ACTIVE"`
	ServerGrpcFoodStoryUrl                        string `env:"SERVER_GRPC_FOOD_STORY_URL"`
	ServerCampaignStartDate                       string `env:"SERVER_CAMPAIGN_START_DATE"`
	ServerCampaignEndDate                         string `env:"SERVER_CAMPAIGN_END_DATE"`
	ServerCampaignActive                          string `env:"SERVER_CAMPAIGN_ACTIVE"`
	ServerCampaignTitle                           string `env:"SERVER_CAMPAIGN_TITLE"`
	ServerSpecialSuperTypes                       string `env:"SERVER_SPECIAL_SUPER_TYPES"`
	ServerFreshPageSuperTypes                     string `env:"SERVER_FRESH_PAGE_SUPER_TYPES"`
	ServerSearchCosmeticId                        string `env:"SERVER_SEARCH_COSMETIC_ID"`
	ServerFoodPartyHost                           string `env:"SERVER_FOOD_PARTY_HOST"`
	ServerApiSecretKey                            string `env:"SERVER_API_SECRET_KEY"`
	ServerRankSecretKey                           string `env:"SERVER_RANK_SECRET_KEY"`
	ServerBiddingDomain                           string `env:"SERVER_BIDDING_DOMAIN"`
	ServerApiDebugKey                             string `env:"SERVER_API_DEBUG_KEY"`
	ServerNewVendorList                           string `env:"SERVER_NEW_VENDOR_LIST"`
	ServerVendorListBannerIgnore                  string `env:"SERVER_VENDOR_LIST_BANNER_IGNORE"`
	ServerVendorListFoodPartyRank                 string `env:"SERVER_VENDOR_LIST_FOOD_PARTY_RANK"`
	ServerFoodpartyCampaignSort                   string `env:"SERVER_FOODPARTY_CAMPAIGN_SORT"`
	ServerMarketQueryLimit                        string `env:"SERVER_MARKET_QUERY_LIMIT"`
	ServerGrpcShopAttributeUrl                    string `env:"SERVER_GRPC_SHOP_ATTRIBUTE_URL"`
	ServerGrpcShopAttributeActive                 string `env:"SERVER_GRPC_SHOP_ATTRIBUTE_ACTIVE"`
	ServerMarketPartyBypass                       string `env:"SERVER_MARKET_PARTY_BYPASS"`
	ZoneInfo                                      string `env:"ZONEINFO"`
	ServerPwaHybridLandingSuperTypes              string `env:"SERVER_PWA_HYBRID_LANDING_SUPER_TYPES"`
	ServerPwaHybridLandingSuperTypeCountThreshold string `env:"SERVER_PWA_HYBRID_LANDING_SUPER_TYPE_COUNT_THRESHOLD"`
	ServerGrpcCpcActive                           string `env:"SERVER_GRPC_CPC_ACTIVE"`
	ServerGrpcCpcUrl                              string `env:"SERVER_GRPC_CPC_URL"`
	ServerExcludeSupermarket                      string `env:"SERVER_EXCLUDE_SUPERMARKET"`
	ServerEtaActive                               string `env:"SERVER_ETA_ACTIVE"`
	ServerNewFilterUpdateChannels                 string `env:"SERVER_NEW_FILTER_UPDATE_CHANNELS"`
	ServerFoodstoryUpdateChannels                 string `env:"SERVER_FOODSTORY_UPDATE_CHANNELS"`
	ServerNewhomeFoodstoryUpdateChannels          string `env:"SERVER_NEWHOME_FOODSTORY_UPDATE_CHANNELS"`
	ServerAllowedOrigins                          string `env:"SERVER_ALLOWED_ORIGINS"`
	ServerNewRelicActive                          string `env:"SERVER_NEW_RELIC_ACTIVE"`
	ServerNewRelicConfigName                      string `env:"SERVER_NEW_RELIC_CONFIG_NAME"`
	ServerNewRelicConfigLicense                   string `env:"SERVER_NEW_RELIC_CONFIG_LICENSE"`
	ServerNewRelicConfigProxy                     string `env:"SERVER_NEW_RELIC_CONFIG_PROXY"`
	ServerTrendingWordsUpdateChannels             string `env:"SERVER_TRENDING_WORDS_UPDATE_CHANNELS"`
	ServerVendorStatus                            string `env:"SERVER_VENDOR_STATUS"`
	ServerServicesUpdateChannels                  string `env:"SERVER_SERVICES_UPDATE_CHANNELS"`
	ServerDistStepSpecificPolygon                 string `env:"SERVER_DIST_STEP_SPECIFIC_POLYGON"`
	ServerPickupOpen                              string `env:"SERVER_PICKUP_OPEN"`
	ServerCategoryOpen                            string `env:"SERVER_CATEGORY_OPEN"`
	ServerSortOpen                                string `env:"SERVER_SORT_OPEN"`
	ServerWebsiteSupertypeExclude                 string `env:"SERVER_WEBSITE_SUPERTYPE_EXCLUDE"`
	ServerNewDesignUpdateChannels                 string `env:"SERVER_NEW_DESIGN_UPDATE_CHANNELS"`
	ServerQuickAccessUpdateChannels               string `env:"SERVER_QUICK_ACCESS_UPDATE_CHANNELS"`
	ServerNewDesignServicesUpdateChannels         string `env:"SERVER_NEW_DESIGN_SERVICES_UPDATE_CHANNELS"`
	ServerSearchiaHost                            string `env:"SERVER_SEARCHIA_HOST"`
	ServerSearchiaIndex                           string `env:"SERVER_SEARCHIA_INDEX"`
	ServerSearchiaApiKey                          string `env:"SERVER_SEARCHIA_API_KEY"`
	ServerGrpcBiddingUpdateChannels               string `env:"SERVER_GRPC_BIDDING_UPDATE_CHANNELS"`
	ServerSpecialSuperTypeNewVersion              string `env:"SERVER_SPECIAL_SUPER_TYPE_NEW_VERSION"`
	ServerCashbackHost                            string `env:"SERVER_CASHBACK_HOST"`
	ServerPanelDbUser                             string `env:"SERVER_PANEL_DB_USER"`
	ServerPanelDbPassword                         string `env:"SERVER_PANEL_DB_PASSWORD"`
	ServerPanelDbHost                             string `env:"SERVER_PANEL_DB_HOST"`
	ServerPanelDbPort                             string `env:"SERVER_PANEL_DB_PORT"`
	ServerPanelDbDb                               string `env:"SERVER_PANEL_DB_DB"`
	ServerFoodProPlanId                           string `env:"SERVER_FOOD_PRO_PLAN_ID"`
	ServerEtaUpdateChannels                       string `env:"SERVER_ETA_UPDATE_CHANNELS"`
	ServerEtaOwnDeliveryUpdateChannels            string `env:"SERVER_ETA_OWN_DELIVERY_UPDATE_CHANNELS"`
	ServerSearchiaUpdateChannels                  string `env:"SERVER_SEARCHIA_UPDATE_CHANNELS"`
	ServerGrpcUbsUrl                              string `env:"SERVER_GRPC_UBS_URL"`
	ServerGrpcUbsActive                           string `env:"SERVER_GRPC_UBS_ACTIVE"`
	ServerProUpdateChannels                       string `env:"SERVER_PRO_UPDATE_CHANNELS"`
	ServerProBannerActive                         string `env:"SERVER_PRO_BANNER_ACTIVE"`
	ServerPartyDisabledOnMenu                     string `env:"SERVER_PARTY_DISABLED_ON_MENU"`

	ServerEcoFoodCollectionId                      string `env:"SERVER_ECO_FOOD_COLLECTION_ID"`
	ServerTopCollectionId                          string `env:"SERVER_TOP_COLLECTION_ID"`
	ServerGrpcFeedbackActive                       string `env:"SERVER_GRPC_FEEDBACK_ACTIVE"`
	ServerGrpcFeedbackUrl                          string `env:"SERVER_GRPC_FEEDBACK_URL"`
	ServerSortBadge                                string `env:"SERVER_SORT_BADGE"`
	ServerUserBaseSortActiveByDefaultChannel       string `env:"SERVER_USER_BASE_SORT_ACTIVE_BY_DEFAULT_CHANNEL"`
	ServerCategoryRenderTypeUpdateChannels         string `env:"SERVER_CATEGORY_RENDER_TYPE_UPDATE_CHANNELS"`
	ServerSortsWithBadge                           string `env:"SERVER_SORTS_WITH_BADGE"`
	ServerIranianEcoTag                            string `env:"SERVER_IRANIAN_ECO_TAG"`
	ServerSearchiaSuggestionIndex                  string `env:"SERVER_SEARCHIA_SUGGESTION_INDEX"`
	ServerSearchiaSuggestionApiKey                 string `env:"SERVER_SEARCHIA_SUGGESTION_API_KEY"`
	ServerGetUserFromHeader                        string `env:"SERVER_GET_USER_FROM_HEADER"`
	ServerProLandingUrl                            string `env:"SERVER_PRO_LANDING_URL"`
	ServerProLandingUrlBonyan                      string `env:"SERVER_PRO_LANDING_URL_BONYAN"`
	ServerProExpirationThreshold                   string `env:"SERVER_PRO_EXPIRATION_THRESHOLD"`
	ServerSearchiaEnablePersonalization            string `env:"SERVER_SEARCHIA_ENABLE_PERSONALIZATION"`
	ServerSuperTypeRankBoosterUpdateChannels       string `env:"SERVER_SUPER_TYPE_RANK_BOOSTER_UPDATE_CHANNELS"`
	ServerCommissionSortUpdateChannels             string `env:"SERVER_COMMISSION_SORT_UPDATE_CHANNELS"`
	ServerNewImpressionUpdateChannels              string `env:"SERVER_NEW_IMPRESSION_UPDATE_CHANNELS"`
	ServerReadFromPanel                            string `env:"SERVER_READ_FROM_PANEL"`
	ServerPickOfTheWeekSortActiveByDefaultChannels string `env:"SERVER_PICK_OF_THE_WEEK_SORT_ACTIVE_BY_DEFAULT_CHANNELS"`
	ServerDisableFilterSortIds                     string `env:"SERVER_DISABLE_FILTER_SORT_IDS"`
	ServerAdsBannerHashSalt                        string `env:"SERVER_ADS_BANNER_HASH_SALT"`
	ServerUbsDismissedCountThreshold               string `env:"SERVER_UBS_DISMISSED_COUNT_THRESHOLD"`
	ServerProMembershipId                          string `env:"SERVER_PRO_MEMBERSHIP_ID"`
	ServerFreeDeliveryMembershipId                 string `env:"SERVER_FREE_DELIVERY_MEMBERSHIP_ID"`

	// Service registry configuration
	ServerServiceRegistryType        string `env:"SERVER_SERVICE_REGISTRY_TYPE"`         // Type of service registry (consul, zookeeper, etcd)
	ServerServiceRegistryAddress     string `env:"SERVER_SERVICE_REGISTRY_ADDRESS"`      // Address of the service registry
	ServerServiceRegistryToken       string `env:"SERVER_SERVICE_REGISTRY_TOKEN"`        // Authentication token for the service registry
	ServerServiceRegistryUsername    string `env:"SERVER_SERVICE_REGISTRY_USERNAME"`     // Username for the service registry
	ServerServiceRegistryPassword    string `env:"SERVER_SERVICE_REGISTRY_PASSWORD"`     // Password for the service registry
	ServerServiceName                string `env:"SERVER_SERVICE_NAME"`                  // Name of this service
	ServerServicePort                string `env:"SERVER_SERVICE_PORT"`                  // Port of this service
	ServerServiceAddress             string `env:"SERVER_SERVICE_ADDRESS"`               // Address of this service
	ServerServiceTags                string `env:"SERVER_SERVICE_TAGS"`                  // Tags for this service (comma-separated)
	ServerServiceHealthCheckURL      string `env:"SERVER_SERVICE_HEALTH_CHECK_URL"`      // Health check URL for this service
	ServerServiceHealthCheckTTL      string `env:"SERVER_SERVICE_HEALTH_CHECK_TTL"`      // Health check TTL for this service
	ServerServiceHealthCheckInterval string `env:"SERVER_SERVICE_HEALTH_CHECK_INTERVAL"` // Health check interval
	ServerServiceScheme              string `env:"SERVER_SERVICE_SCHEME"`                // Service URL scheme (http/https)
}

func (c Config) GetString(key string) string {
	// Convert the key to lowercase for case-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a case-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-insensitive)
		if strings.ToLower(envTag) == lowerKey {
			if field.IsValid() && field.Kind() == reflect.String {
				return field.String()
			}
			break
		}
	}

	// Fallback to exact tag name
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-sensitive)
		if envTag == key {
			if field.IsValid() && field.Kind() == reflect.String {
				return field.String()
			}
			break
		}
	}

	return ""
}

func (c Config) GetInt(key string) int {
	// Convert the key to lowercase for case-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a case-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-insensitive)
		if strings.ToLower(envTag) == lowerKey {
			if field.IsValid() && field.Kind() == reflect.String {
				if intValue, err := strconv.Atoi(field.String()); err == nil {
					return intValue
				}
			}
			break
		}
	}

	// Fallback to exact tag name
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-sensitive)
		if envTag == key {
			if field.IsValid() && field.Kind() == reflect.String {
				if intValue, err := strconv.Atoi(field.String()); err == nil {
					return intValue
				}
			}
			break
		}
	}

	return 0
}

func (c Config) GetBool(key string) bool {
	// Convert the key to lowercase for case-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a case-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-insensitive)
		if strings.ToLower(envTag) == lowerKey {
			if field.IsValid() && field.Kind() == reflect.String {
				if boolValue, err := strconv.ParseBool(field.String()); err == nil {
					return boolValue
				}
			}
			break
		}
	}

	// Fallback to exact tag name
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-sensitive)
		if envTag == key {
			if field.IsValid() && field.Kind() == reflect.String {
				if boolValue, err := strconv.ParseBool(field.String()); err == nil {
					return boolValue
				}
			}
			break
		}
	}

	return false
}

func (c Config) GetStringSlice(key string) []string {
	// Convert the key to lowercase for case-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a case-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-insensitive)
		if strings.ToLower(envTag) == lowerKey {
			if field.IsValid() && field.Kind() == reflect.String {
				return []string{field.String()}
			}
			break
		}
	}

	// Fallback to exact tag name
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-sensitive)
		if envTag == key {
			if field.IsValid() && field.Kind() == reflect.String {
				return []string{field.String()}
			}
			break
		}
	}

	return nil
}
