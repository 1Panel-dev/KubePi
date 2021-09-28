let App = angular.module('App', []);

App.controller('IndexCtrl', function ($scope, $http, $log) {

    const PREFIX = "@WK@_";
    const TYPE_KUBE_CONFIG = "Kube Config";
    const TYPE_TOKEN = "Token";

    $scope.reset = function () {
        $("#file").val("");
        $(".custom-file-label").text("Upload a kube config");
        $scope.msg = "";
        $scope.item = {
            type: TYPE_KUBE_CONFIG
        };
    };

    $("#file").bind('change', function (e) {
        let f = e.target.files[0];
        let reader = new FileReader();
        reader.readAsText(f || "");
        reader.onload = function () {
            $(".custom-file-label").text(f.name);
            $scope.item.kubeConfig = window.btoa(unescape(encodeURIComponent(this.result)));
        }
    });

    $scope.list = function () {
        $scope.items = [];
        for (let i = 0; i < localStorage.length; i++) {
            let key = localStorage.key(i);
            if (key.includes(PREFIX)) {
                $scope.items.push(angular.fromJson(localStorage.getItem(key)));
            }
        }
    };

    $scope.connect = function (item) {
        const CONFIG_URL = "api/kube-config";
        const TOKEN_URL = "api/kube-token";
        let url = item.type === TYPE_KUBE_CONFIG ? CONFIG_URL : TOKEN_URL;
        return $http.post(url, item).then(function (response) {
            if (response.data.success) {
                let shellUrl = location.pathname + "terminal?token=" + response.data.token;
                window.open(shellUrl, "_blank");
            } else {
                $scope.error(response.data.message);
                $log.error(response.data.message);
            }
        }, function (response) {
            if (response.data.message) {
                $scope.error(response.data.message);
            } else {
                $scope.error(response);
            }
            $log.error(response)
        });
    };

    $scope.confirm = function (item) {
        $scope.deleteItem = item;
    };

    $scope.delete = function () {
        localStorage.removeItem($scope.deleteItem.key);
        $scope.deleteItem = null;
        $scope.list();
    };

    $scope.error = function (msg) {
        if (angular.isObject(msg)) {
            msg = angular.toJson(msg);
        }
        let html =
            '<div class="alert alert-danger alert-dismissible fade show mt-3" role="alert">\n' +
            '                ' + msg +
            '                <button type="button" class="close" data-dismiss="alert" aria-label="Close">\n' +
            '                    <span aria-hidden="true">&times;</span>\n' +
            '                </button>\n' +
            '            </div>';
        $(".clearfix").after(html);
    };

    $scope.save = function () {
        if (!$scope.item.name) {
            $scope.msg = "Name is required.";
            return;
        }
        let item = localStorage.getItem(PREFIX + $scope.item.name);
        if (item) {
            $scope.msg = "Name already exist.";
            return;
        }

        if ($scope.item.type === TYPE_KUBE_CONFIG) {
            if (!$scope.item.kubeConfig) {
                $scope.msg = "Kube config is required.";
                return;
            }
        }

        if ($scope.item.type === TYPE_TOKEN) {
            if (!$scope.item.apiServer) {
                $scope.msg = "ApiServer is required.";
                return;
            }
            if (!$scope.item.token) {
                $scope.msg = "Token is required.";
                return;
            }
        }

        $scope.item.key = PREFIX + $scope.item.name;
        localStorage.setItem($scope.item.key, angular.toJson($scope.item));

        // clear
        $scope.reset();
        $scope.list();
        $('#add').modal('hide');
    };

    $scope.reset();
    $scope.list();
});