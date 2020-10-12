import React from "react"
import Layout from "../components/layout"

interface IndexPageProps {}

const IndexPage: React.FC<IndexPageProps> = () => (
  <Layout>
    <h1>Hi people</h1>
    <p>Welcome to your new Gatsby site.</p>
    <p>Now go build something great.</p>
  </Layout>
)

export default IndexPage
