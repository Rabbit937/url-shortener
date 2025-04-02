import express from 'express'
import cors from 'cors'
const app = express()

app.set('view engine', 'ejs');
app.use(express.static('public'));
app.use(cors())

app.get('/', (req, res) => {
  res.render('index');
});

app.listen(3000, () => {
  console.log('Node server is running on http://localhost:3000');
})