'use strict';

angular.module('app.main', ['ngRoute'])

    .config(['$routeProvider', function($routeProvider) {
        $routeProvider.when('/main', {
            templateUrl: 'views/main/main.html',
            controller: 'MainCtrl'
        });
    }])

    .controller('MainCtrl', function($scope, $rootScope, $http, toastr) {

        $rootScope.server = JSON.parse(localStorage.getItem("old_darkID_server"));

        $scope.generatingID = false;
        $scope.ids = [];
        $http.get(clientapi + 'ids')
            .then(function(data) {
                console.log('data success');
                console.log(data);
                $scope.ids = data.data;

            }, function(data) {
                console.log('data error');
            });

        $scope.newID = function() {
            $scope.generatingID = true;
            $http.get(clientapi + 'newid')
                .then(function(data) {
                    console.log('data success');
                    console.log(data);
                    $scope.ids = data.data;
                    $scope.generatingID = false;

                }, function(data) {
                    console.log('data error');
                });
        };

        $scope.blindAndSendToSign = function(id) {
            $http.get(clientapi + 'blindandsendtosign/' + id)
                .then(function(data) {
                    console.log('data success');
                    console.log(data);
                    $scope.ids = data.data;

                }, function(data) {
                    console.log('data error');
                });
        };
        $scope.verify = function(id) {
            $http.get(clientapi + 'verify/' + id)
                .then(function(data) {
                    console.log('data success');
                    console.log(data);
                    $scope.ids = data.data;

                }, function(data) {
                    console.log('data error');
                });
        };
        $scope.clientApp = function(route, param) {
            $http.get(clientapi + route + '/' + param)
                .then(function(data) {
                    console.log('data success');
                    console.log(data);
                    $scope.ids = data.data;

                }, function(data) {
                    console.log('data error');
                });
        };
        $scope.addToBlockchain = function(id) {
            $http({
                    url: 'http://127.0.0.1:3002/register',
                    method: "POST",
                    headers: {
                        "Content-Type": undefined
                    },
                    data: {address: id}
                })
                .then(function(data) {
                    //$scope.ids = data.data;
                    toastr.success("added to blockchain");
                },
                function(data) {
                    console.log(data);
                });
        };
    });
