'use strict';

angular.module('eventprint', [
  'ngResource',
  'ngRoute',
  'mm.foundation'
])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });
  });
