## 产品

- product
- product_translation: 产品拆分表

- commodity_versions: 产品版本表
    - product_id: 一个产品有多个版本
    - status: new、configuring、configured

    - commodity_version_service_type: 产品版本和服务配置关联表
        - service_blueprint_id:一个产品版本有多个服务配置
        - commodity_version_id:

        - commodity_version_service_type_parameter_group:
            - commodity_version_service_type_id:

- 层级关系
    - 产品 ===> 产品版本 ===> 服务配置 ===> 参数组
    - 一个产品有多个产品版本
    - 一个产品版本有多个服务配置
    - 一个服务配置对应一个参数组

## 项目

- project
- project2product: 项目和产品的关联表

## 服务

- service: 关联 project_id、service_blueprint_id

## 服务配置

- service_blueprint: 关联 product_id、resource_id
- service_blueprint_translation: 服务配置拆分表

## 资源包

- resource_packs: 关联 service_id

## 关系

- 一个产品中有多个项目、一个产品中可以有多个服务配置、一个产品有多个产品版本、一个产品版本有多个服务配置
- 一个项目中有多个服务、
- 一个服务中有多个资源包、一个服务只属于一个项目
- 一个服务配置可以被多个服务使用

## 举例

- 产品(项目1、项目2、服务配置1、服务配置2)
- 项目1(服务1、服务1)
- 服务1(资源包1、资源包2、资源包3、服务配置1)
- 独立: 资源包、服务配置