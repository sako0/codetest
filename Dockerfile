FROM ruby:3.0.0
RUN apt-get update -qq && apt-get install -y nodejs
WORKDIR /app
COPY Gemfile /app/Gemfile
COPY Gemfile.lock /app/Gemfile.lock
RUN bundle install
COPY . /app
COPY entrypoint.sh /usr/bin/
RUN chmod +x /usr/bin/entrypoint.sh
ENTRYPOINT ["entrypoint.sh"]

EXPOSE 8888
CMD ["rails", "server", "-b", "0.0.0.0","-p","8888"]