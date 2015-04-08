// Ionic Starter App

// angular.module is a global place for creating, registering and retrieving Angular modules
// 'starter' is the name of this angular module example (also set in a <body> attribute in index.html)
// the 2nd parameter is an array of 'requires'
// 'starter.services' is found in services.js
// 'starter.controllers' is found in controllers.js
angular.module('starter', ['ionic', 'starter.controllers', 'starter.services'])

    .run(function ($ionicPlatform, $rootScope, $state, auth) {
        $ionicPlatform.ready(function () {
            // Hide the accessory bar by default (remove this to show the accessory bar above the keyboard
            // for form inputs)
            if (window.cordova && window.cordova.plugins && window.cordova.plugins.Keyboard) {
                cordova.plugins.Keyboard.hideKeyboardAccessoryBar(true);
            }
            if (window.StatusBar) {
                // org.apache.cordova.statusbar required
                StatusBar.styleLightContent();
            }
        });

        $rootScope.$on('$stateChangeStart', function (event, toState, toParams) {
            var requireLogin = toState.data.requireLogin;
            console.log(toState);
            if (requireLogin && typeof $rootScope.currentUser === 'undefined') {
                event.preventDefault();
                return $state.go('login');
            }
        });
    })

    .config(function ($stateProvider, $urlRouterProvider, $httpProvider) {

        $httpProvider.interceptors.push(function ($timeout, $q, $injector) {
            var auth, $http, $state;
            // this trick must be done so that we don't receive
            // `Uncaught Error: [$injector:cdep] Circular dependency found`
            $timeout(function () {
                auth = $injector.get('auth');
                $http = $injector.get('$http');
                $state = $injector.get('$state');
            });

            return {
                response: function (rejection) {
                    var deferred = $q.defer();
                    if (rejection.status == 200) {
                        if(rejection.data && rejection.data.Type == "login" && rejection.data.Status ==401 && rejection.config.url!= window.localStorage['loginurl']){
                            auth.logindlg()
                                .then(function () {
                                    deferred.resolve();
                                    $state.go(toState.name, toParams);
                                }).catch(function () {
                                    deferred.reject(rejection);
                                    $state.go('login');
                                });
                        }else{
                            deferred.resolve(rejection);
                        }
                    }else{
                        deferred.resolve(rejection);
                    }
                    return deferred.promise;
                }
            };
        });
        // Ionic uses AngularUI Router which uses the concept of states
        // Learn more here: https://github.com/angular-ui/ui-router
        // Set up the various states which the app can be in.
        // Each state's controller can be found in controllers.js
        $stateProvider
            // setup an abstract state for the tabs directive
            .state('login', {
                url: "/login",
                //abstract: true,
                templateUrl: "tpl/login.html",
                controller: 'LoginCtrl',
                data: {
                    requireLogin: false
                }
            })
            .state('tab', {
                url: "/tab",
                abstract: true,
                templateUrl: "tpl/tabs.html",
                data: {
                    requireLogin: false
                }
            })

            // Each tab has its own nav history stack:

            .state('tab.dash', {
                url: '/dash',
                views: {
                    'tab-dash': {
                        templateUrl: 'tpl/tab-dash.html',
                        controller: 'DashCtrl'
                    }
                },
                data: {
                    requireLogin: true
                }
            })

            .state('tab.chats', {
                url: '/chats',
                views: {
                    'tab-chats': {
                        templateUrl: 'tpl/tab-chats.html',
                        controller: 'ChatsCtrl'
                    }
                },
                data: {
                    requireLogin: true
                }
            })
            .state('tab.chat-detail', {
                url: '/chats/:chatId',
                views: {
                    'tab-chats': {
                        templateUrl: 'tpl/chat-detail.html',
                        controller: 'ChatDetailCtrl'
                    }
                },
                data: {
                    requireLogin: true
                }
            })

            .state('tab.account', {
                url: '/account',
                views: {
                    'tab-account': {
                        templateUrl: 'tpl/tab-account.html',
                        controller: 'AccountCtrl'
                    }
                },
                data: {
                    requireLogin: true
                }
            });

        // if none of the above states are matched, use this as the fallback
        $urlRouterProvider.otherwise('/login');

    });