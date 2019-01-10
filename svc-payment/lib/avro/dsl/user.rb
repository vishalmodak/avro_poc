record :user do
  required :id, :int
  required :name, :string
  required :email, :string
  # required :location, :string, default: 'chicago'
end
