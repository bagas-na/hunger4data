package grpc

// type subscriptionHandler struct {
// 	subscriptionclient pb.Subscription_ServiceClient
// }

// func NewHandAuth(subscriptionclient pb.Subscription_ServiceClient) *subscriptionHandler {
// 	return &subscriptionHandler{
// 		subscriptionclient: subscriptionclient,
// 	}
// }

// func (h *subscriptionHandler) GetCountries(c echo.Context) error {
// 	var req model.Country
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{
// 			"error": "invalid request body",
// 		})
// 	}
// 	ctx := context.TODO()
// 	resp, err := h.subscriptionclient.Get_Countries(ctx, &pb.Empty{})
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"message": resp.Message,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"data": resp.Countries,
// 	})
// }

// func (h *subscriptionHandler) CreateSub(c echo.Context) error {
// 	var req model.Subscription
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{
// 			"error": "invalid request body",
// 		})
// 	}
// 	ctx := context.TODO()
// 	resp, err := h.subscriptionclient.Create_Subscription(ctx, &pb.Subscription_Request{
// 		IdUser:    req.Id_user.String(),
// 		IdCountry: req.Id_country.String(),
// 	})
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"message": resp.Message,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"message": resp.Message,
// 	})
// }

// func (h *subscriptionHandler) GetSubByID(c echo.Context) error {
// 	var req model.Subscription
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{
// 			"error": "invalid request body",
// 		})
// 	}
// 	ctx := context.TODO()
// 	resp, err := h.subscriptionclient.Get_Subscription_By_ID(ctx, &pb.Subscription_Request{
// 		Id:        req.Id.String(),
// 		IdUser:    req.Id_user.String(),
// 		IdCountry: req.Id_country.String(),
// 	})

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"message": resp.Message,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"data": resp.Subscription,
// 	})
// }

// func (h *subscriptionHandler) UpdateSub(c echo.Context) error {
// 	var req model.Subscription
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{
// 			"error": "invalid request body",
// 		})
// 	}
// 	ctx := context.TODO()
// 	resp, err := h.subscriptionclient.Update_Subscription(ctx, &pb.Subscription_Request{
// 		Id:        req.Id.String(),
// 		IdUser:    req.Id_user.String(),
// 		IdCountry: req.Id_country.String(),
// 	})
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"message": resp.Message,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"message": resp.Message,
// 	})
// }

// func (h *subscriptionHandler) DeleteSub(c echo.Context) error {
// 	var req model.Subscription
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{
// 			"error": "invalid request body",
// 		})
// 	}
// 	ctx := context.TODO()
// 	resp, err := h.subscriptionclient.Delete_Subscription(ctx, &pb.Subscription_Request{
// 		Id:        req.Id.String(),
// 		IdUser:    req.Id_user.String(),
// 		IdCountry: req.Id_country.String(),
// 	})
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"message": resp.Message,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"message": resp.Message,
// 	})
// }
