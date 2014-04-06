'use strict';

angular.module('eventprint')
  .controller('MainCtrl', function ($scope, $filter, Attendee) {
    $scope.filteredAttendees = [];
    $scope.attendees = Attendee.query(function () {
      $scope.filterAttendees();
    });
    $scope.currentPage = 1;
    $scope.pageSize = 25;
    $scope.maxPaginationButtons = 7;

    $scope.filterAttendees = function () {
      var begin = (($scope.currentPage - 1) * $scope.pageSize),
          end = begin + $scope.pageSize;

      $scope.filteredAttendees = $filter('filter')($scope.attendees, $scope.query);
      $scope.filteredAttendees = $filter('orderBy')($scope.filteredAttendees, 'lastName');
      $scope.filteredAttendees = $scope.filteredAttendees.slice(begin, end);
    };

    $scope.checkIn = function (attendee) {
      attendee.$checkIn();
    };

    $scope.$watch('currentPage + pageSize', $scope.filterAttendees);
    $scope.$watch('query', $scope.filterAttendees);
  });
