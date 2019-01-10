Rails.application.routes.draw do
  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
  namespace :v1 do
    get '/payments/:loan_number' => "payment#lookup_payments_loan_number"
  end
end
