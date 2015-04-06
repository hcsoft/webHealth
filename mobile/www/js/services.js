angular.module('starter.services', ['ionic'])

    .factory('Chats', function () {
        // Might use a resource here that returns a JSON array

        // Some fake testing data
        var chats = [{
            id: 0,
            name: 'Ben Sparrow',
            lastText: 'You on your way?',
            face: 'https://pbs.twimg.com/profile_images/514549811765211136/9SgAuHeY.png'
        }, {
            id: 1,
            name: 'Max Lynx',
            lastText: 'Hey, it\'s me',
            face: 'https://avatars3.githubusercontent.com/u/11214?v=3&s=460'
        }, {
            id: 2,
            name: 'Andrew Jostlin',
            lastText: 'Did you get the ice cream?',
            face: 'https://pbs.twimg.com/profile_images/491274378181488640/Tti0fFVJ.jpeg'
        }, {
            id: 3,
            name: 'Adam Bradleyson',
            lastText: 'I should buy a boat',
            face: 'https://pbs.twimg.com/profile_images/479090794058379264/84TKj_qa.jpeg'
        }, {
            id: 4,
            name: 'Perry Governor',
            lastText: 'Look at my mukluks!',
            face: 'https://pbs.twimg.com/profile_images/491995398135767040/ie2Z_V6e.jpeg'
        }];

        return {
            all: function () {
                return chats;
            },
            remove: function (chat) {
                chats.splice(chats.indexOf(chat), 1);
            },
            get: function (chatId) {
                for (var i = 0; i < chats.length; i++) {
                    if (chats[i].id === parseInt(chatId)) {
                        return chats[i];
                    }
                }
                return null;
            }
        };
    }).service('auth', function ($ionicModal, $rootScope, $q, $state, $http) {
        function assignCurrentUser(_username) {
            $rootScope.currentUser = _username;
        }

        return {
            login: function (_username, _password) {
                var deferred = $q.defer();
                //进行http请求查询结果
                console.log(_username, _password);
                $http.jsonp("http://localhost:3000/login?callback=JSON_CALLBACK",
                    {
                        params: {
                            userid: _username,
                            password: CryptoJS.MD5(_password).toString()
                        }
                    }).then(function (data) {
                        console.log(data);
                        if (data.data.Status == 200) {
                            deferred.resolve(data);
                        } else {
                            deferred.reject(data);
                        }
                    }).catch(function (data) {
                        console.log(data);
                        deferred.reject(data);
                    });
                //deferred.resolve();
                return deferred.promise.then(assignCurrentUser);
            },
            logindlg: function () {
                var deferred = $q.defer();
                var child = $rootScope.$new();
                child.cancel = function () {
                    //取消登录时,啥也不干
                    this.modal.hide();
                    deferred.reject();
                };
                child.login = function (_username, _password) {
                    this.modal.hide();
                    deferred.resolve(_username);
                };
                var instance = $ionicModal.fromTemplateUrl('tpl/logindialog.html', {
                    scope: child,
                    animation: 'slide-in-up',
                    focusFirstInput: true,
                    backdropClickToClose: false
                }).then(function (modal) {
                    child.modal = modal;
                    modal.show();
                });

                return deferred.promise.then(assignCurrentUser);
            },
            logout: function () {
                $rootScope.currentUser = undefined;
                $state.go("login")
            }
        };

    });
