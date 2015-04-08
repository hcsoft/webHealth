angular.module('starter.controllers', [])

    .controller('DashCtrl', function ($scope, $http, $timeout) {
        $scope.querydata = {
            querystr: ''
        }
        $scope.query = function (event) {
            if (!event || event.keyCode == 13) {
                //清空列表
                $scope.listdata = [];

                querydata("name")
                querydata("idnumber")
                querydata("address")

            }
            return true;
        };
        function querydata(type) {
            $timeout(
                function () {
                    $http.jsonp(window.localStorage['remoteurl'] + "admin/query/file" + window.localStorage['callback'],
                        {
                            params: {
                                district: $scope.$parent.currentdistid,
                                querystring: $scope.querydata.querystr,
                                querytype: type
                            }
                        }
                    ).then(function (data) {
                            console.log(data);
                            if(data.data)
                                $scope.listdata = $scope.listdata.concat(data.data);
                        });
                }
            )
        }

        $scope.$parent.$watch("currentdistid", function () {
            $scope.query();
        });

    }).controller('IndexCtrl', function ($scope, $http) {
        $scope.currentdistid = null;
        $http.jsonp(window.localStorage['remoteurl'] + "admin/district" + window.localStorage['callback']).then(function (data) {
            console.log(data);
            $scope.dist = {child: [data.data]};
            $scope.$parent.currentdistid = data.data['ID'];
            $scope.$parent.currentdistname = data.data['Name'];
        });
        if (!window.localStorage['remoteurl'])
            window.localStorage['remoteurl'] = 'http://192.168.168.2:3000/';
        window.localStorage['callback'] = '?callback=JSON_CALLBACK';
        window.localStorage['loginurl'] = window.localStorage['remoteurl'] + 'login' + window.localStorage['callback']
        $scope.reloadDist = function () {
            console.log("reloadDist");
            if (!$scope.dist) {
                $http.jsonp(window.localStorage['remoteurl'] + "admin/district" + window.localStorage['callback']).then(function (data) {
                    $scope.dist = {child: [data.data]};
                    $scope.$parent.currentdistid = data.data['ID'];
                    $scope.$parent.currentdistname = data.data['Name'];

                });
            }
        }
    })

    .controller('LoginCtrl', function ($scope, auth, $state) {
        $scope.login = function (uid, pwd) {
            auth.login(uid, pwd).then(function () {
                $state.go("tab.dash")
            }).catch(function (data) {
                if (data.status == 404) {
                    $scope.msg = "登录超时";
                } else if (data.data && data.data.Msg) {
                    $scope.msg = data.data.Msg;
                } else {
                    $scope.msg = "登录失败!";
                }
            });
        }
    })
    .controller('ChatsCtrl', function ($scope, Chats) {
        $scope.chats = Chats.all();
        $scope.remove = function (chat) {
            Chats.remove(chat);
        }
    })
    .controller('ContentController', function ($scope, $ionicSideMenuDelegate) {
        $scope.toggleLeft = function () {
            $ionicSideMenuDelegate.toggleLeft();
        }
    })
    .controller('ChatDetailCtrl', function ($scope, $stateParams, Chats) {
        $scope.chat = Chats.get($stateParams.chatId);
    })

    .controller('AccountCtrl', function ($scope, auth) {
        $scope.settings = {
            removete: true
        };
        $scope.logout = function () {
            auth.logout();
        };
        $scope.remoteurl = window.localStorage['remoteurl'];
        $scope.setSettings = function (remoteurl) {
            $scope.remoteurl = remoteurl
            window.localStorage['remoteurl'] = $scope.remoteurl;
            window.localStorage['callback'] = '?callback=JSON_CALLBACK';
            window.localStorage['loginurl'] = window.localStorage['remoteurl'] + 'login' + window.localStorage['callback']
        }

    }).controller('LoginModalCtrl', function ($scope, $state) {


    });
