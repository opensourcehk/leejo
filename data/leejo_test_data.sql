--
-- These are dummy data for integration test
--

insert into leejo_user (username, password, gender)
  values ('test', 'test', 'F');

insert into leejo_api_client (id, secret, redirect_uri)
  values ('testing', 'testing', 'http://localhost:8000/redirect');
