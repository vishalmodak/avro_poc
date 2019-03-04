FROM ruby:2.4-slim

RUN apt-get update && apt-get install -y -f build-essential nodejs libpq-dev protobuf-compiler iputils-ping

#copy protos from domain_protos images
RUN mkdir -p /domain-protos/build/distributions/
COPY --from=domain_protos /home/gradle/app/build/distributions/domain-protos.tgz /domain-protos/build/distributions/

RUN gem install bundler -v 2.0.1

RUN mkdir /app
WORKDIR /app

COPY Gemfile Gemfile.lock ./
RUN bundle install

COPY . .

EXPOSE 3000

#CMD puma -C config/puma.rb
CMD ["bundle", "exec", "rails", "server", "-b", "0.0.0.0"]