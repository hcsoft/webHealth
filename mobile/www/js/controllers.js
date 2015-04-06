angular.module('starter.controllers', [])

    .controller('DashCtrl', function ($scope) {
    })
    .controller('LoginCtrl', function ($scope, auth , $state) {
        $scope.login = function(uid,pwd){
            auth.login(uid,pwd).then(function(){
                $state.go("tab.dash")
            }).catch(function (data) {
                console.log(data);
                $scope.msg=data.data.Msg;
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

    .controller('AccountCtrl', function ($scope,auth) {
        $scope.settings = {
            enableFriends: true
        };
        $scope.logout = function(){
            auth.logout();
        }
    }).controller('LoginModalCtrl', function ($scope,$state) {


    });
