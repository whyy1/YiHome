function getCookie(name) {
    var r = document.cookie.match("\\b" + name + "=([^;]*)\\b");
    return r ? r[1] : undefined;
}

// 解析提取url中的查询字符串参数
function decodeQuery(){
    var search = decodeURI(document.location.search);
    return search.replace(/(^\?)/, '').split('&').reduce(function(result, item){
        values = item.split('=');
        result[values[0]] = values[1];
        return result;
    }, {});
}
$(document).ready(function(){

    $.get("/api/v1.0/areas", function (resp) {
        if ("0" == resp.errno) {
            // // 表示查询到了数据,修改前端页面
            // for (var i=0; i<resp.data.length; i++) {
            //     // 向页面中追加标签
            //     var areaId = resp.data[i].aid;
            //     var areaName = resp.data[i].aname;
            //     $("#area-id").append('<option value="'+ areaId +'">'+ areaName +'</option>');
            // }

            // 使用前端模板
            rendered_html = template("areas-tmpl", {areas: resp.data});
            $("#area-id").html(rendered_html);
        } else {
            alert(resp.errmsg);
        }
    }, "json");

    // 获取详情页面要展示的房屋编号
    var queryData = decodeQuery();
    var houseId = queryData["id"];
    console.log(houseId);

    // 获取该房屋的详细信息
    $.get("/api/v1.0/houses/" + houseId, function(resp){
        if ("0" == resp.errno) {
            // $(".swiper-container").html(template("house-image-tmpl", {img_urls:resp.data.house.img_urls, price:resp.data.house.price}));
            // $(".detail-con").html(template("house-detail-tmpl", {house:resp.data.house}));
            console.log(resp.data);
            $("#house-title").val(resp.data.house.title);
            $("#house-price").val(resp.data.house.deposit);
            $("#house-address").val(resp.data.house.address);
            $("#house-room-count").val(resp.data.house.room_count);
            $("#house-acreage").val(resp.data.house.acreage);
            $("#house-unit").val(resp.data.house.unit);
            $("#house-capacity").val(resp.data.house.capacity);
            $("#house-beds").val(resp.data.house.beds);
            $("#house-deposit").val(resp.data.house.deposit);
            $("#house-min-days").val(resp.data.house.min_days);
            $("#house-max-days").val(resp.data.house.max_days);
            if(!resp.data.house.max_days){
                $("#house-max-days").val(0);
            }
            $(".house-facility-list clearfix").html(resp.data.house.facilities);
            let inputlables =  $("[name=facility]")
            for (const label of inputlables.toArray()) {
               if(resp.data.house.facilities.find((e)=> e == parseInt(label.value))) {
                label.checked = true
               }
            }
            
        }
    })

    // 处理房屋基本信息的表单数据
    $("#form-house-info").submit(function (e) {
        e.preventDefault();
        // 检验表单数据是否完整
        // 将表单的数据形成json，向后端发送请求
        var formData = {};
        $(this).serializeArray().map(function (x) { formData[x.name] = x.value });

        // 对于房屋设施的checkbox需要特殊处理
        var facility = [];
        // $("input:checkbox:checked[name=facility]").each(function(i, x){ facility[i]=x.value });
        $(":checked[name=facility]").each(function(i, x){ facility[i]=x.value });

        formData.facility = facility;

        //发送请求时显示加载动画效果
        var loadingimg = document.getElementById('loading');
        loadingimg.style.display = 'block';
    
        // 使用ajax向后端发送请求
        $.ajax({
            url: "/api/v1.0/houses/"+houseId,
            type: "put",
            data: JSON.stringify(formData),
            contentType: "application/json",
            dataType: "json",
            timeout:10000,
            headers: {
                "X-CSRFToken": getCookie("csrf_token")
            },
            success: function(resp){
                if ("4101" == resp.errno) {
                    location.href = "/home/login.html";
                } else if ("0" == resp.errno) {
                    alert('修改成功！');
                    location.href = "/home/myhouse.html";
                } else {
                    alert(resp.errmsg);
                }
                loadingimg.style.display = 'none';
            },
            error: function(){
                alert("网络异常，请求接收超时");
                loadingimg.style.display = 'none';
            }
        })
    })
})