'use strict';

angular.module('eventprint')
  .controller('MainCtrl', function ($scope, $filter, $debounce, Attendee) {
    $scope.paginatedAttendees = [];
    $scope.filteredAttendees = [];
    $scope.attendees = Attendee.query(function () {
      $scope.filterAndPaginateAttendees();
    });
    $scope.alerts = [];
    $scope.currentPage = 1;
    $scope.pageSize = 25;
    $scope.maxPaginationButtons = 7;

    $scope.filterAndPaginateAttendees = function () {
      var begin = (($scope.currentPage - 1) * $scope.pageSize),
          end = begin + $scope.pageSize;

      $scope.filteredAttendees = $filter('filter')($scope.attendees, $scope.query);
      $scope.filteredAttendees = $filter('orderBy')($scope.filteredAttendees, 'lastName');
      $scope.paginatedAttendees = $scope.filteredAttendees.slice(begin, end);
    };

    $scope.sync = function () {
      $scope.addAlert('info', 'Sync started. We\'ll report back to you when it\'s done.');

      Attendee.sync(function () {
        $scope.addAlert('success', 'Sync compeleted successfully. Refresh at will.');
      });
    }

    $scope.checkIn = function (attendee) {
      attendee.$checkIn();
    };

    $scope.addAlert = function (type, msg) {
      $scope.alerts.push({type: type, msg: msg});
    };

    $scope.closeAlert = function (index) {
      $scope.alerts.splice(index, 1);
    };

    $scope.$watch('currentPage + pageSize', $scope.filterAndPaginateAttendees);
    $scope.$watch('query', function (newValue, oldValue) {
      if (newValue === oldValue) {
        return;
      }

      $debounce(applyQuery, 350, false);
    });

    var applyQuery = function () {
      $scope.currentPage = 1;
      $scope.filterAndPaginateAttendees();
    }
  });
