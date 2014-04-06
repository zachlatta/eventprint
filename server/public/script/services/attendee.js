angular.module('eventprint')
  .factory('Attendee', function ($resource) {
    return $resource('/api/attendees/:id/:action',
      {
        id: '@id',
        action: '@action'
      },
      {
        checkIn: {
          method: 'PUT',
          params: {
            action: 'check_in'
          }
        }
      });
  })
