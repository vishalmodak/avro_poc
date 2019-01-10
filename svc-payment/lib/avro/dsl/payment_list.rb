record :payment_list do
  required :payments, :array, items: :payment
end
