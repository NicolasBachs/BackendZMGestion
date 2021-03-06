package controllers

import (
	"BackendZMGestion/internal/db"
	"BackendZMGestion/internal/gestores"
	"BackendZMGestion/internal/helpers"
	"BackendZMGestion/internal/interfaces"
	"BackendZMGestion/internal/models"
	"BackendZMGestion/internal/structs"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mitchellh/mapstructure"
)

//RolesController contiene los metodos para los endpoints de la API referidos a Roles.
type RolesController struct {
	DbHandler *db.DbHandler
}

/**
 * @api {POST} /roles/dame Dame Rol
 * @apiPermission Administradores
 * @apiDescription Devuelve un rol a partir de un Id
 * @apiGroup Roles
 * @apiParam {Object} Roles
 * @apiParam {int} Roles.IdRol
 * @apiParamExample {json} Request-Example:
 {
	 "Roles": {
		 "IdRol":2
	 }
 }
 * @apiSuccessExample {json} Success-Response:
 {
    "error": null,
    "respuesta": {
		"Roles": {
            "IdRol": 2,
            "Rol": "Vendedores",
            "FechaAlta": "2020-04-09 15:01:35.000000",
            "Descripcion": "Este rol es para los vendedores"
        }
	}
}
* @apiErrorExample {json} Error-Response:
 {
    "error": {
        "codigo": "ERROR_NOEXISTE_ROL",
        "mensaje": "No existe el rol."
    },
    "respuesta": null
}
*/
//Dame Devuelve un rol a partir de un Id
func (rc *RolesController) Dame(c echo.Context) error {

	rol := structs.Roles{}

	jsonMap, err := helpers.GenerateMapFromContext(c)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
	}

	mapstructure.Decode(jsonMap["Roles"], &rol)

	//_ = c.Request().Header.Get("Authorization")
	rolesService := models.RolesService{
		DbHandler: rc.DbHandler,
		Rol:       &rol,
	}
	result, err := rolesService.Dame()

	if err != nil || result == nil {
		return interfaces.GenerarRespuestaError(err, http.StatusBadRequest)
	}

	response := interfaces.Response{
		Error: nil,
	}

	response.AddModels(result)

	return c.JSON(http.StatusOK, response)
}

/**
 * @api {POST} /roles/crear Crear Rol
 * @apiPermission Administradores
 * @apiDescription Permite crear un rol
 * @apiGroup Roles
 * @apiHeader {String} Authorization
 * @apiParam {Object} Roles
 * @apiParam {string} Roles.Rol
 * @apiParam {string} [Roles.Descripcion]
 * @apiParamExample {json} Request-Example:
 {
	 "Roles": {
		 "Rol": "Encargados"
	 }
 }
 * @apiSuccessExample {json} Success-Response:
 {
    "error": null,
    "respuesta": {
		"Roles" : {
			"IdRol": 7,
			"Rol": "Encargados",
			"FechaAlta": "2020-04-09 15:01:35.000000",
			"Descripcion": ""
		}
	}
}
* @apiErrorExample {json} Error-Response:
 {
    "error": {
        "codigo": "ERROR_EXISTE_NOMBREROL",
        "mensaje": "El nombre de rol ya existe."
    },
    "respuesta": null
}
*/
func (rc *RolesController) Crear(c echo.Context) error {
	rol := structs.Roles{}

	jsonMap, err := helpers.GenerateMapFromContext(c)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
	}

	mapstructure.Decode(jsonMap["Roles"], &rol)

	headerToken := c.Request().Header.Get("Authorization")
	token, err := helpers.GetToken(headerToken)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
	}

	gestorRoles := gestores.GestorRoles{
		DbHandler: rc.DbHandler,
	}

	result, err := gestorRoles.Crear(rol, *token)

	if err != nil || result == nil {
		return interfaces.GenerarRespuestaError(err, http.StatusBadRequest)
	}

	response := interfaces.Response{
		Error: nil,
	}

	response.AddModels(result)

	return c.JSON(http.StatusOK, response)
}

/**
 * @api {GET} /roles Listar Roles
 * @apiPermission Administradores
 * @apiDescription Devuelve una lista de roles
 * @apiGroup Roles
 * @apiSuccessExample {json} Success-Response:
 {
    "error": null,
    "respuesta": [
		{
			"Roles":{
				"IdRol": 1,
				"Rol": "Administradores",
				"FechaAlta": "2020-04-09 15:01:35.000000",
				"Observaciones": ""
			}
		},
		{
			"Roles":{
				"IdRol": 2,
				"Rol": "Vendedores",
				"FechaAlta": "2020-04-09 15:01:35.000000",
				"Observaciones": ""
			}
		}

    ]
}
* @apiErrorExample {json} Error-Response:
 {
    "error": {
        "codigo": "ERROR_DEFAULT",
        "mensaje": "Ha ocurrido un error mientras se procesaba su petición."
    },
    "respuesta": null
}
*/
//Listar Lista los roles
func (rc *RolesController) Listar(c echo.Context) error {

	gestorRoles := gestores.GestorRoles{
		DbHandler: rc.DbHandler,
	}

	result, err := gestorRoles.Listar()

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusBadRequest)
	}

	response := interfaces.Response{
		Error: nil,
	}

	var respuesta []map[string]interface{}
	for _, el := range result {
		objeto := make(map[string]interface{})
		objeto["Roles"] = el
		respuesta = append(respuesta, objeto)
	}
	response.Respuesta = respuesta

	return c.JSON(http.StatusOK, response)
}

/**
 * @api {POST} /roles/borrar Borrar Rol
 * @apiPermission Administradores
 * @apiDescription Borra un rol a partir de su Id
 * @apiGroup Roles
 * @apiHeader {String} Authorization
 * @apiParam {Object} Roles
 * @apiParam {int} Roles.IdRol
 * @apiParamExample {json} Request-Example:
 {
	 "Roles": {
		 "IdRol": "2"
	 }
 }
 * @apiSuccessExample {json} Success-Response:
 {
    "error": null,
    "respuesta": null
}
*/
//Borrar Devuelve un rol a partir de un Id
func (rc *RolesController) Borrar(c echo.Context) error {

	rol := structs.Roles{}

	jsonMap, err := helpers.GenerateMapFromContext(c)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
	}

	mapstructure.Decode(jsonMap["Roles"], &rol)

	//_ = c.Request().Header.Get("Authorization")
	gestorRoles := gestores.GestorRoles{
		DbHandler: rc.DbHandler,
	}

	headerToken := c.Request().Header.Get("Authorization")
	token, err := helpers.GetToken(headerToken)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
	}

	err = gestorRoles.Borrar(rol, *token)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusBadRequest)
	}

	response := interfaces.Response{
		Error:     nil,
		Respuesta: nil,
	}

	return c.JSON(http.StatusOK, response)
}

/**
 * @api {POST} /roles/modificar Modificar Rol
 * @apiPermission Administradores
 * @apiDescription Permite modificar un rol
 * @apiGroup Roles
 * @apiHeader {String} Authorization
 * @apiParam {Object} Roles
 * @apiParam {int} Roles.IdRol
 * @apiParam {string} Roles.Rol
 * @apiParam {string} [Roles.Descripcion]
 * @apiParamExample {json} Request-Example:
 {
	 "Roles": {
		 "IdRol": 7,
		 "Rol": "Los encargados",
		 "Descripcion": "Nueva descripcion"
	 }
 }
 * @apiSuccessExample {json} Success-Response:
 {
    "error": null,
    "respuesta": {
		"IdRol": 7,
		"Rol": "Los encargados",
		"FechaAlta": "2020-04-09 15:01:35.000000",
		"Descripcion": "Nueva descripcion"
	}
}
* @apiErrorExample {json} Error-Response:
 {
    "error": {
        "codigo": "ERROR_EXISTE_NOMBREROL",
        "mensaje": "El nombre de rol ya existe."
    },
    "respuesta": null
}
*/
func (rc *RolesController) Modificar(c echo.Context) error {
	rol := structs.Roles{}

	jsonMap, err := helpers.GenerateMapFromContext(c)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
	}

	mapstructure.Decode(jsonMap["Roles"], &rol)

	headerToken := c.Request().Header.Get("Authorization")
	token, err := helpers.GetToken(headerToken)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
	}

	gestorRoles := gestores.GestorRoles{
		DbHandler: rc.DbHandler,
	}

	result, err := gestorRoles.Modificar(rol, *token)

	if err != nil || result == nil {
		return interfaces.GenerarRespuestaError(err, http.StatusBadRequest)
	}

	response := interfaces.Response{
		Error: nil,
	}

	response.AddModels(result)

	return c.JSON(http.StatusOK, response)
}

/**
 * @api {POST} /roles/listarPermisos Listar Permisos
 * @apiPermission Administradores
 * @apiName Listar Permisos
 * @apiDescription Devuelve una lista de permisos para un determinado rol
 * @apiGroup Roles
 * @apiParam {Object} Roles
 * @apiParam {int} IdRol
 * @apiSuccessExample {json} Success-Response:
 {
    "error": null,
    "respuesta": [
		{
            "Permisos": {
                "IdPermiso": 1,
                "Permiso": "Crear rol"
            }
        },
        {
            "Permisos": {
                "IdPermiso": 2,
                "Permiso": "Borrar rol"
            }
        }
    ]
}
* @apiErrorExample {json} Error-Response:
 {
    "error": {
        "codigo": "ERROR_DEFAULT",
        "mensaje": "Ha ocurrido un error mientras se procesaba su petición."
    },
    "respuesta": null
}
*/
//Listar Lista los permisos para un rol
func (rc *RolesController) ListarPermisos(c echo.Context) error {

	rol := structs.Roles{}

	jsonMap, err := helpers.GenerateMapFromContext(c)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
	}

	mapstructure.Decode(jsonMap["Roles"], &rol)

	rolesService := models.RolesService{
		DbHandler: rc.DbHandler,
		Rol:       &rol,
	}

	result, err := rolesService.ListarPermisos()

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusBadRequest)
	}

	response := interfaces.Response{
		Error: nil,
	}

	var respuesta []map[string]interface{}
	for _, el := range result {
		objeto := make(map[string]interface{})
		objeto["Permisos"] = el
		respuesta = append(respuesta, objeto)
	}
	response.Respuesta = respuesta

	return c.JSON(http.StatusOK, response)
}

/**
 * @api {POST} /roles/asignarPermisos Asignar Permisos
 * @apiPermission Administradores
 * @apiDescription Permite asignar los permisos a un determinado rol
 * @apiGroup Roles
 * @apiParam {Object} Roles
 * @apiParam {int} Roles.IdRol
 * @apiParam {Object[]} Permisos
 * @apiParamExample {json} Request-Example:
 {
	"Roles": {
		"IdRol": 9
	},
	"Permisos": [
		{
			"IdPermiso": 3
		},
		{
			"IdPermiso": 4
		}
	]
}
 * @apiSuccessExample {json} Success-Response:
 {
    "error": null,
    "respuesta": null
}
* @apiErrorExample {json} Error-Response:
 {
    "error": {
        "codigo": "ERROR_DEFAULT",
        "mensaje": "Ha ocurrido un error mientras se procesaba su petición."
    },
    "respuesta": null
}
*/
func (rc *RolesController) AsignarPermisos(c echo.Context) error {
	var permisos []structs.Permisos
	var body Request

	err := json.NewDecoder(c.Request().Body).Decode(&body)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
	}

	for _, el := range body.Permisos {
		var permiso structs.Permisos
		err = mapstructure.Decode(el, &permiso)
		if err != nil {
			return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
		}
		permisos = append(permisos, permiso)
	}

	rolesService := models.RolesService{
		DbHandler: rc.DbHandler,
		Rol:       body.Roles,
	}

	headerToken := c.Request().Header.Get("Authorization")
	token, err := helpers.GetToken(headerToken)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusUnprocessableEntity)
	}

	err = rolesService.AsignarPermisos(permisos, *token)

	if err != nil {
		return interfaces.GenerarRespuestaError(err, http.StatusBadRequest)
	}

	response := interfaces.Response{
		Error:     nil,
		Respuesta: nil,
	}

	return c.JSON(http.StatusOK, response)
}

type Request struct {
	Permisos []structs.Permisos `json:"Permisos"`
	Roles    *structs.Roles     `json:"Roles"`
}
