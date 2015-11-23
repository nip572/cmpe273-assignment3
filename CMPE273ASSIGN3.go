package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "encoding/json"
    "strings"
    "io/ioutil"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "strconv"
    "sort"
    "bytes"
    "log"
)



//ESTIMATE PRICE STRUCTURE
type Estimate_Price struct {
    StartLatitude  float64
    StartLongitude float64
    EndLatitude    float64
    EndLongitude   float64
    Prices         []EstimatesPrice `json:"prices"`
}


type EstimatesPrice struct {
    ProductId       string  `json:"product_id"`
    CurrencyCode    string  `json:"currency_code"`
    DisplayName     string  `json:"display_name"`
    Estimate        string  `json:"estimate"`
    LowEstimate     int     `json:"low_estimate"`
    HighEstimate    int     `json:"high_estimate"`
    SurgeMultiplier float64 `json:"surge_multiplier"`
    Duration        int     `json:"duration"`
    Distance        float64 `json:"distance"`
}



const (

 API_URL string = "https://sandbox-api.uber.com/v1/%s%s"  // URL FOR SANDBOX

)

func (pe *Estimate_Price) get(c *Client) error {
    EstimatesPriceParams := map[string]string{
        "start_latitude":  strconv.FormatFloat(pe.StartLatitude, 'f', 2, 32),
        "start_longitude": strconv.FormatFloat(pe.StartLongitude, 'f', 2, 32),
        "end_latitude":    strconv.FormatFloat(pe.EndLatitude, 'f', 2, 32),
        "end_longitude":   strconv.FormatFloat(pe.EndLongitude, 'f', 2, 32),
    }

    data := c.getRequest("estimates/price", EstimatesPriceParams)
    if e := json.Unmarshal(data, &pe); e != nil {
        return e
    }
    return nil
}

//GETTER INTERFACE
type Getter interface {

    get(c *Client) error
}

type RequestOptions struct {
    ServerToken    string
    ClientId       string
    ClientSecret   string
    AppName        string
    AuthorizeUrl   string
    AccessTokenUrl string
    AccessToken string
    BaseUrl        string
}


type Client struct {

    Options *RequestOptions
}


func Create(options *RequestOptions) *Client {

    return &Client{options}

}


func (c *Client) Get(getter Getter) error {

    if e := getter.get(c); e != nil {
        return e
    }

    return nil
}
//GET REQUEST
func (c *Client) getRequest(endpoint string, params map[string]string) []byte {
    urlParams := "?"
    params["server_token"] = c.Options.ServerToken
    for k, v := range params {
        if len(urlParams) > 1 {
            urlParams += "&"
        }
        urlParams += fmt.Sprintf("%s=%s", k, v)
    }

    url := fmt.Sprintf(API_URL, endpoint, urlParams)

    res, err := http.Get(url)
    if err != nil {
        
    }

    data, err := ioutil.ReadAll(res.Body)
    res.Body.Close()

    return data
}
//END GET REQUEST

type Products struct {
    Latitude  float64
    Longitude float64
    Products  []Product `json:"products"`
}


type Product struct {
    ProductId   string `json:"product_id"`
    Description string `json:"description"`
    DisplayName string `json:"display_name"`
    Capacity    int    `json:"capacity"`
    Image       string `json:"image"`
}


func (pl *Products) get(c *Client) error {
    productParams := map[string]string{
        "latitude":  strconv.FormatFloat(pl.Latitude, 'f', 2, 32),
        "longitude": strconv.FormatFloat(pl.Longitude, 'f', 2, 32),
    }

    data := c.getRequest("products", productParams)
    if e := json.Unmarshal(data, &pl); e != nil {
        return e
    }
    return nil
}



type request_object struct{
Id int
Name string `json:"Name"`
Address string `json:"Address"`
City string `json:"City"`
State string `json:"State"`
Zip string `json:"Zip"`
Coordinates struct{
    Lat float64
    Lng float64
}
}

var id int;
var tripId int;


type Responz struct {
    Results []struct {
        AddressComponents []struct {
            LongName  string   `json:"long_name"`
            ShortName string   `json:"short_name"`
            Types     []string `json:"types"`
        } `json:"address_components"`
        FormattedAddress string `json:"formatted_address"`
        Geometry         struct {
            Location struct {
                Lat float64 `json:"lat"`
                Lng float64 `json:"lng"`
            } `json:"location"`
            LocationType string `json:"location_type"`
            Viewport     struct {
                Northeast struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"northeast"`
                Southwest struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"southwest"`
            } `json:"viewport"`
        } `json:"geometry"`
        PartialMatch bool     `json:"partial_match"`
        PlaceID      string   `json:"place_id"`
        Types        []string `json:"types"`
    } `json:"results"`
    Status string `json:"status"`
}

type TripResponse struct {    
    ID                     string   `json:"id"`
    Status                 string   `json:"status"`    
    StartingFromLocationID string   `json:"starting_from_location_id"`
    BestRouteLocationIds   []string `json:"best_route_location_ids"`
    TotalUberCosts         int      `json:"total_uber_costs"`
    TotalUberDuration      int      `json:"total_uber_duration"`
    TotalDistance          float64  `json:"total_distance"`
}

type RideRequest struct {
    EndLatitude    string `json:"end_latitude"`
    EndLongitude   string `json:"end_longitude"`
    ProductID      string `json:"product_id"`
    StartLatitude  string `json:"start_latitude"`
    StartLongitude string `json:"start_longitude"`
}

type OnGoingTrip struct {
    BestRouteLocationIds      []string `json:"best_route_location_ids"`
    ID                        string   `json:"id"`
    NextDestinationLocationID string   `json:"next_destination_location_id"`
    StartingFromLocationID    string   `json:"starting_from_location_id"`
    Status                    string   `json:"status"`
    TotalDistance             float64  `json:"total_distance"`
    TotalUberCosts            int      `json:"total_uber_costs"`
    TotalUberDuration         int      `json:"total_uber_duration"`
    UberWaitTimeEta           int      `json:"uber_wait_time_eta"`
}

type ReqResponse struct {

    Driver          interface{} `json:"driver"`
    Eta             int         `json:"eta"`
    Location        interface{} `json:"location"`
    RequestID       string      `json:"request_id"`
    Status          string      `json:"status"`
    SurgeMultiplier int         `json:"surge_multiplier"`
    Vehicle         interface{} `json:"vehicle"`

}




type resObj struct{
Greeting string
}
//LOCATION CREATE     
func LOCATION_CREATE(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    id=id+1;


    decoder := json.NewDecoder(req.Body)
    var t request_object 
    t.Id = id; 
    err := decoder.Decode(&t)
    if err != nil {
        fmt.Println("Error")
    }


    
    st:=strings.Join(strings.Split(t.Address," "),"+");
    fmt.Println(st);
    constr := []string {strings.Join(strings.Split(t.Address," "),"+"),strings.Join(strings.Split(t.City," "),"+"),t.State}
    lstringplus := strings.Join(constr,"+")
    locstr := []string{"http://maps.google.com/maps/api/geocode/json?address=",lstringplus}

    resp, err := http.Get(strings.Join(locstr,""))

    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       fmt.Println("Error: Wrong address");
     }
     var data Responz
    err = json.Unmarshal(body, &data)
    fmt.Println(data.Status)
  
    t.Coordinates.Lat=data.Results[0].Geometry.Location.Lat;
    t.Coordinates.Lng=data.Results[0].Geometry.Location.Lng;




 conn, err := mgo.Dial("mongodb://nipun:nipun@ds045464.mongolab.com:45464/db2")

  
    if err != nil {
        panic(err)
    }
    defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("db2").C("qwerty");
err = c.Insert(t);


    js,err := json.Marshal(t)
    if err != nil{
	   fmt.Println("Error")
	   return
	}
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}
//END CREATE LOCATION

//START GET LOCATION

func GET_LOCATION(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
fmt.Println(p.ByName("locid"));
id ,err1:= strconv.Atoi(p.ByName("locid"))
if err1 != nil {
        panic(err1)
    }
 conn, err := mgo.Dial("mongodb://nipun:nipun@ds045464.mongolab.com:45464/db2")


    if err != nil {
        panic(err)
    }
    defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("db2").C("qwerty");
result:=request_object{}
err = c.Find(bson.M{"id":id}).One(&result)
if err != nil {
                fmt.Println(err)
        }

    
        js,err := json.Marshal(result)
    if err != nil{
       fmt.Println("Error")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}

type modrequest_object struct{
    Address string `json:"address"`
    City string `json:"city"`
    State string `json:"state"`
    Zip string `json:"zip"`
}


//END GET LOCATION

//START UPDATE LOCATION
func UPDATE_LOCATION(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
   
 id ,err1:= strconv.Atoi(p.ByName("locid"))

 if err1 != nil {
         panic(err1)
     }
  conn, err := mgo.Dial("mongodb://nipun:nipun@ds045464.mongolab.com:45464/db2")

     if err != nil {
         panic(err)
     }
     defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
 c:=conn.DB("db2").C("qwerty");


     decoder := json.NewDecoder(req.Body)
     var t modrequest_object  
     err = decoder.Decode(&t)
     if err != nil {
         fmt.Println("Error")
     }


     colQuerier := bson.M{"id": id}
     change := bson.M{"$set": bson.M{"address": t.Address, "city":t.City,"state":t.State,"zip":t.Zip}}
     err = c.Update(colQuerier, change)
     if err != nil {
         panic(err)
     }

}
//END UPDATE LOCATION

//START DELETE LOCATION
func DELETE_LOCATION(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
     id ,err1:= strconv.Atoi(p.ByName("locid"))

 if err1 != nil {
         panic(err1)
     }
  conn, err := mgo.Dial("mongodb://nipun:nipun@ds045464.mongolab.com:45464/db2")
  conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("db2").C("qwerty");


     if err != nil {
         panic(err)
     }
     defer conn.Close();
     err=c.Remove(bson.M{"id":id})
     if err != nil { fmt.Printf("Could not find kitten %s to delete", id)}
    rw.WriteHeader(http.StatusNoContent)
}

type userUber struct {
    LocationIds            []string `json:"location_ids"`
    StartingFromLocationID string   `json:"starting_from_location_id"`
}

//END DELETE LOCATION

func TRIP_PLANNING(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

    decoder := json.NewDecoder(req.Body)
    var uUD userUber 
    err := decoder.Decode(&uUD)
    if err != nil {
        log.Println("Error")
    }

        log.Println(uUD.StartingFromLocationID);


//UBER DETAILS
    var options RequestOptions;

    options.ServerToken= "ZNSdsdziWkOSLcdrrsy0RhMiZI7_i_GlGzpA3f1L";
    
    options.ClientId= "5-BNiHDpt1CZvQoWd2G2vV2GSvSnIu2j";
    
    options.ClientSecret= "P5qyGJI-sJw5m-s2kFHljzg59kccexZ8qkbaL44P";
    
    options.AppName= "";

    options.BaseUrl= "https://sandbox-api.uber.com/v1/";
    

    client :=Create(&options); 


        sid ,err1:= strconv.Atoi(uUD.StartingFromLocationID)
 fmt.Println(uUD.StartingFromLocationID);
 fmt.Println(sid);
 if err1 != nil {
         panic(err1)
     }

    conn, err := mgo.Dial("mongodb://nipun:nipun@ds045464.mongolab.com:45464/db2");

    if err != nil {
        panic(err)
    }
    defer conn.Close();

    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("db2").C("qwerty");
    result:=request_object{}
    err = c.Find(bson.M{"id":sid}).One(&result)
    if err != nil {
                fmt.Println(err)
        }

 
    index:=0;
    totalPrice := 0;
    totalDistance :=0.0;
    totalDuration :=0;
    bestroute:=make([]float64,len(uUD.LocationIds));
    m := make(map[float64]string)

    for _,ids := range uUD.LocationIds{
    
        lid,err1:= strconv.Atoi(ids)
           
        if err1 != nil {
            panic(err1)
        }
        

        resultLID:=request_object{}
        err = c.Find(bson.M{"id":lid}).One(&resultLID)
        if err != nil {
             fmt.Println(err)
        }
        pe := &Estimate_Price{}
        pe.StartLatitude = result.Coordinates.Lat;

        pe.StartLongitude = result.Coordinates.Lng;

        pe.EndLatitude = resultLID.Coordinates.Lat;

        pe.EndLongitude = resultLID.Coordinates.Lng;

        if e := client.Get(pe); e != nil {
            fmt.Println(e);
        }
        totalDistance=totalDistance+pe.Prices[0].Distance;

        totalDuration=totalDuration+pe.Prices[0].Duration;

        totalPrice=totalPrice+pe.Prices[0].LowEstimate;

        bestroute[index]=pe.Prices[0].Distance;

        m[pe.Prices[0].Distance]=ids;
        index=index+1;
    }
  
    sort.Float64s(bestroute);
    


   

    var tripres TripResponse;

    tripId=tripId+1;

     tripres.ID=strconv.Itoa(tripId);

     tripres.TotalDistance=totalDistance;

     tripres.TotalUberCosts=totalPrice;

     tripres.TotalUberDuration=totalDuration;

     tripres.Status="Planning";

     tripres.StartingFromLocationID=strconv.Itoa(sid);

     tripres.BestRouteLocationIds=make([]string,len(uUD.LocationIds));

     index=0;

     for _, ind := range bestroute{
        tripres.BestRouteLocationIds[index]=m[ind];
        index=index+1;
     }
     fmt.Println(tripres.BestRouteLocationIds[1]);



    c1:=conn.DB("db2").C("trips");

    err = c1.Insert(tripres);


        js,err := json.Marshal(tripres)
    if err != nil{
       fmt.Println("Error")

       return
    }
    rw.Header().Set("Content-Type","application/json")

    rw.Write(js)

    }


func GET_TRIP(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

    conn, err := mgo.Dial("mongodb://nipun:nipun@ds045464.mongolab.com:45464/db2")

    
    if err != nil {

        panic(err)
    }
    defer conn.Close();

    conn.SetMode(mgo.Monotonic,true);

    c:=conn.DB("db2").C("trips");

    result:=TripResponse{}

    err = c.Find(bson.M{"id":p.ByName("tripid")}).One(&result)

    if err != nil {

        fmt.Println(err)
    }


    js,err := json.Marshal(result)

    if err != nil{

       fmt.Println("Error")

       return
    }
    rw.Header().Set("Content-Type","application/json")

    rw.Write(js)
}


var currentPos int;

var ogtID int;

func REQUEST_TRIP(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    


   
    kid ,err1:= strconv.Atoi(p.ByName("tripid"))
    var siD int;
    
    if err1 != nil {
         panic(err1)
     }
    var ogt OnGoingTrip;

    result1:=request_object{}

    result2:=request_object{}

    conn, err := mgo.Dial("mongodb://nipun:nipun@ds045464.mongolab.com:45464/db2")

   
    if err != nil {

        panic(err)
    }
    defer conn.Close();

    conn.SetMode(mgo.Monotonic,true);

    c:=conn.DB("db2").C("trips");

    result:=TripResponse{}

    err = c.Find(bson.M{"id":strconv.Itoa(kid)}).One(&result)

    if err != nil {

        fmt.Println(err)
    }else{

    var iD int;

    c1:=conn.DB("db2").C("qwerty");

    if currentPos==0{

        iD, err = strconv.Atoi(result.StartingFromLocationID)

        siD=iD;

        if err != nil {
     
            fmt.Println(err)
        }
    }else
    {
        iD, err = strconv.Atoi(result.BestRouteLocationIds[currentPos-1])
        siD=iD;
        if err != nil {
      
            fmt.Println(err)
        }
    }

    err = c1.Find(bson.M{"id":iD}).One(&result1)
    if err != nil {
                fmt.Println(err)
        }
    iD, err = strconv.Atoi(result.BestRouteLocationIds[currentPos])
    if err != nil {
        
        fmt.Println(err)
    }
    err = c1.Find(bson.M{"id":iD}).One(&result2)
    if err != nil {
                fmt.Println(err)
        }


        fmt.Println(result2.Coordinates.Lat);
    }

    ogt.ID=strconv.Itoa(ogtID);

    ogt.BestRouteLocationIds=result.BestRouteLocationIds;

    ogt.StartingFromLocationID=strconv.Itoa(siD);

    ogt.NextDestinationLocationID=result.BestRouteLocationIds[currentPos];

    ogt.TotalDistance=result.TotalDistance;

    ogt.TotalUberCosts=result.TotalUberCosts;

    ogt.TotalUberDuration=result.TotalUberDuration;

    ogt.Status="requesting";
    
    var options RequestOptions;

    options.ServerToken= "ZNSdsdziWkOSLcdrrsy0RhMiZI7_i_GlGzpA3f1L";

    options.ClientId= "5-BNiHDpt1CZvQoWd2G2vV2GSvSnIu2j";

    options.ClientSecret= "P5qyGJI-sJw5m-s2kFHljzg59kccexZ8qkbaL44P";

    options.AppName= "";

    options.BaseUrl= "https://sandbox-api.uber.com/v1/";

    client :=Create(&options);

    pl:=Products{};

    pl.Latitude=result1.Coordinates.Lat;

    pl.Longitude=result1.Coordinates.Lng;

    if e := pl.get(client); e != nil {

         fmt.Println(e)
    }

    var prodid string;

    i:=0

    for _, product := range pl.Products {

         if(i == 0){
             prodid = product.ProductId
        }
    }



    var rr RideRequest;

    rr.StartLatitude=strconv.FormatFloat(result1.Coordinates.Lat, 'f', 6, 64);

    rr.StartLongitude=strconv.FormatFloat(result1.Coordinates.Lng, 'f', 6, 64);

    rr.EndLatitude=strconv.FormatFloat(result2.Coordinates.Lat, 'f', 6, 64);

    rr.EndLongitude=strconv.FormatFloat(result2.Coordinates.Lng, 'f', 6, 64);

    rr.ProductID=prodid;

    buf, _ := json.Marshal(rr)

    body := bytes.NewBuffer(buf)

    url := fmt.Sprintf(API_URL, "requests?","access_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiIyZTBiOGFiOC05MmVkLTQzMmItYjcwNC05ZGRkMmQxY2VlMTMiLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6IjYwNTE1YWM5LTgwNzAtNGExMi1iYmM3LTRlOTQyZTViZDY5YiIsImV4cCI6MTQ1MDc2MTU4MCwiaWF0IjoxNDQ4MTY5NTgwLCJ1YWN0IjoiMjJBVUVkbWhEUmVoU3hMUG5DZGlJMXNtbFNLWEtzIiwibmJmIjoxNDQ4MTY5NDkwLCJhdWQiOiJEX21qMDg5NDZxaERhVG1wUDBZRkJrWHFOdE95ZGJ4RSJ9.U3aoSDQV2kty0Njsb0rMZg4_xWdPRZEae5gq3d7NHnV7Nmv4-gsKTgJrWUDNpwbDdMVo4rbcdJ2dTBLVUs0Wf1vF02DgYbwUjYPlfJuT0ahPJxqOSIim4iGvmOGAWn9NdFIaup7VZyVoj-RlD8bofJ6mlb3YQifrLX8iHZ686-rOkWPz8XaeWsiojENowdrJlh_FESWoE_d3mX48mhMTSt_jCfm5bBZ0bQ4kvVRXP_3ZadaJPiRUMWKKJHYLaI0-j1uplpsScBqgx34IHafvGPwvxrhQ_gltBEDE-xKGvs75m4GdqOGpUBwWjQa5tUFzEwnD9HJzz5HnGR4GbeA4tQ")

    res, err := http.Post(url,"application/json",body)

    if err != nil {

        fmt.Println(err)
    }
    data, err := ioutil.ReadAll(res.Body)

    var rRes ReqResponse;

    err = json.Unmarshal(data, &rRes)

    ogt.UberWaitTimeEta=rRes.Eta;

    js,err := json.Marshal(ogt)

    if err != nil{

       fmt.Println("Error")

       return
    }
    ogtID=ogtID+1;

    currentPos=currentPos+1;

    rw.Header().Set("Content-Type","application/json")

    rw.Write(js)

}


func main() {
    A := httprouter.New()


    id=0;

    tripId=0;

    currentPos=0;

    ogtID=0;
    
    fmt.Println("LISTENING ON PORT 9000")

    A.POST("/locations",LOCATION_CREATE)

    A.POST("/trips",TRIP_PLANNING)

    A.GET("/locations/:locid",GET_LOCATION)

    A.GET("/trips/:tripid",GET_TRIP)

    A.PUT("/locations/:locid",UPDATE_LOCATION)

    A.PUT("/trips/:tripid/request",REQUEST_TRIP)

    A.DELETE("/locations/:locid",DELETE_LOCATION)

    server := http.Server{
            Addr:        "0.0.0.0:9000",
            Handler: A,
    }

    server.ListenAndServe()
}